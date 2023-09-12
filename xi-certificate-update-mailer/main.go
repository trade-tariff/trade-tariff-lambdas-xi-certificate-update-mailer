package main

import (
	"os"
	"time"

	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/emailer"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/fetcher"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/logger"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/parser"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/joho/godotenv"
)

var defaultConfig = map[string]string{
	"PERSISTENCE_BUCKET": "trade-tariff-persistence-844815912454",
	"BUCKET_PREFIX":      "data/taric",
	"DEBUG":              "false",
	"FROM_EMAIL":         "Online Trade Tariff Support <trade-tariff-support@enginegroup.com>",
	"TO_EMAILS":          "William Fish <william.fish@digital.hmrc.gov.uk>",
}

type LambdaEvent struct {
	Date string `json:"date"`
}

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_VERSION") != "" {
		logger.Log.Info("Running in AWS Lambda environment")
		lambda.Start(execute)
	} else {
		logger.Log.Info("Running in local environment")
		execute(nil)
	}
}

func execute(event *LambdaEvent) {
	initializeEnvironment()

	sess := initializeAWSSession()

	s3Client := s3.New(sess)
	fetcher := fetcher.NewFetcher(
		s3Client,
		os.Getenv("PERSISTENCE_BUCKET"),
		os.Getenv("BUCKET_PREFIX"),
	)

	sesClient := ses.New(sess)
	emailer := emailer.NewEmailer(sesClient)

	date := getDateArgument(event)
	processCertificateUpdate(fetcher, emailer, date)
}

func initializeEnvironment() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()

		if err != nil {
			logger.Log.Fatal("Error loading .env file")
		}
	}

	for key, defaultValue := range defaultConfig {
		if val, exists := os.LookupEnv(key); !exists || val == "" {
			err := os.Setenv(key, defaultValue)

			if err != nil {
				logger.Log.Fatal("Error occurred while setting environment variable")
			}
		}
	}
}

func initializeAWSSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		logger.Log.Fatal("Error occurred while creating AWS session. Have you configured your AWS credentials?")
	}
	return sess
}

func getDateArgument(event *LambdaEvent) string {
	if event != nil && event.Date != "" {
		return event.Date
	}

	if len(os.Args) < 2 {
		return time.Now().Format("2006-01-02")
	} else {
		return os.Args[1]
	}
}

func processCertificateUpdate(fetcher *fetcher.Fetcher, e emailer.Emailer, date string) {
	object := fetcher.FetchFileObject(date)
	file := fetcher.FetchXML(object)

	certificateUpdate := parser.ParseXML(file)

	if len(certificateUpdate.CertificateChanges) > 0 {
		emailConfiguration := emailer.NewEmailConfiguration(certificateUpdate)
		e.Send(emailConfiguration)
	} else {
		logger.Log.Info(
			"No certificate changes found",
			logger.String("date", date),
			logger.String("filename", file.Filename()),
			logger.Int("contentLength", len(file.Content)),
		)
	}
}
