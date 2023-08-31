package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type EmailConfiguration struct {
	fromEmail string
	subject   string
	message   string
	toEmails  []*string
	charset   string
}

func fetchMessage(update CertificateUpdate) string {
	filename := os.Getenv("MESSAGE_TEMPLATE_FILENAME")
	file, err := os.ReadFile(filename)

	fmt.Println(filename)
	fmt.Println(string(file))

	for _, change := range update.CertificateChanges {
		fmt.Println(change)
		fmt.Println(change.ValidityStartDate())
	}

	if err != nil {
		log.Fatal("Error opening email template")
	}

	messageTemplate, err := template.New("message").Parse(string(file))

	if err != nil {
		log.Fatalf("Error parsing email template: %s", err)
	}

	writer := new(bytes.Buffer)

	err = messageTemplate.ExecuteTemplate(writer, "message", update)

	if err != nil {
		log.Fatalf("Error executing email template: %s", err)
	}

	return writer.String()
}

func fetchFromEmail() string {
	fromEmail := os.Getenv("FROM_EMAIL")

	if fromEmail == "" {
		fromEmail = "Online Trade Tariff Support <trade-tariff-support@enginegroup.com>"
	}

	return fromEmail
}

func fetchToEmails() []*string {
	toEmailsStr := os.Getenv("TO_EMAILS")

	if toEmailsStr == "" {
		log.Fatal("No TO_EMAILS specified")
	}

	toEmailsList := strings.Split(toEmailsStr, ",")

	var toEmailsPtrs []*string

	for _, email := range toEmailsList {
		toEmailsPtrs = append(toEmailsPtrs, &email)
	}

	return toEmailsPtrs
}

func fetchSubject(filename string) string {
	writer := new(bytes.Buffer)
	subject := os.Getenv("SUBJECT")
	subjectTemplate, err := template.New("subject").Parse(subject)

	if err != nil {
		log.Fatal("Error parsing subject template")
	}

	subjectTemplate.ExecuteTemplate(writer, "subject", map[string]interface{}{
		"Filename": filename,
	})

	return writer.String()
}

func fetchEmailConfiguration(update CertificateUpdate) EmailConfiguration {
	subject := fetchSubject("test.txt")
	message := fetchMessage(update)
	fromEmail := fetchFromEmail()
	toEmails := fetchToEmails()

	return EmailConfiguration{
		fromEmail: fromEmail,
		subject:   subject,
		message:   message,
		toEmails:  toEmails,
	}
}

func SendEmailSES(emailConfig EmailConfiguration) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)
	if err != nil {
		log.Println("Error occurred while creating aws session", err)
		return
	}

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Source: aws.String(emailConfig.fromEmail),
		Destination: &ses.Destination{
			ToAddresses: emailConfig.toEmails,
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String(emailConfig.charset),
				Data:    aws.String(emailConfig.subject),
			},

			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(emailConfig.charset),
					Data:    aws.String(emailConfig.message),
				},
			},
		},
	}

	_, err = svc.SendEmail(input)

	if err != nil {
		log.Println("Error sending mail - ", err)
		return
	}

	var emails []string

	for _, email := range emailConfig.toEmails {
		emails = append(emails, *email)
	}

	log.Println("Email sent successfully to: ", strings.Join(emails, ","))
}
