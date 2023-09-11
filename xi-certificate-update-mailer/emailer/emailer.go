package emailer

import (
	"bytes"
	_ "embed"
	"os"
	"strings"
	"text/template"

	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/logger"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/parser"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

//go:embed templates/subject_template.txt
var subjectTemplate string

//go:embed templates/email_template.html
var emailTemplate string

type EmailConfiguration struct {
	fromEmail string
	subject   string
	body      string
	toEmails  []*string
	charset   string
}

func NewEmailConfiguration(update parser.CertificateUpdate) EmailConfiguration {
	subject := fetchSubject(update)
	message := fetchBody(update)
	fromEmail := fetchFromEmail()
	toEmails := fetchToEmails()

	return EmailConfiguration{
		fromEmail: fromEmail,
		subject:   subject,
		body:      message,
		toEmails:  toEmails,
		charset:   "UTF-8",
	}
}

type SesClient interface {
	SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error)
}

type Emailer struct {
	SesClient SesClient
	Config    EmailConfiguration
}

func NewEmailer(sesClient SesClient) Emailer {
	return Emailer{
		SesClient: sesClient,
	}
}

func (e *Emailer) Send(config EmailConfiguration) {
	input := &ses.SendEmailInput{
		Source: aws.String(config.fromEmail),
		Destination: &ses.Destination{
			ToAddresses: config.toEmails,
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String(config.charset),
				Data:    aws.String(config.subject),
			},

			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(config.charset),
					Data:    aws.String(config.body),
				},
			},
		},
	}

	_, err := e.SesClient.SendEmail(input)

	if err != nil {
		logger.Log.Fatal("Error sending mail")
		return
	}

	var emails []string

	for _, email := range config.toEmails {
		emails = append(emails, *email)
	}

	logger.Log.Info(
		"Email sent successfully",
		logger.String("to", strings.Join(emails, ",")),
		logger.String("subject", config.subject),
	)
}

func fetchBody(update parser.CertificateUpdate) string {
	messageTemplate, err := template.New("message").Parse(emailTemplate)

	if err != nil {
		logger.Log.Fatal("Error parsing email template")
	}

	writer := new(bytes.Buffer)

	err = messageTemplate.ExecuteTemplate(writer, "message", update)

	if err != nil {
		logger.Log.Fatal("Error executing email template")
	}

	return writer.String()
}

func fetchFromEmail() string {
	fromEmail := os.Getenv("FROM_EMAIL")

	if fromEmail == "" {
		logger.Log.Fatal("No FROM_EMAIL specified")
	}

	return fromEmail
}

func fetchToEmails() []*string {
	toEmailsStr := os.Getenv("TO_EMAILS")

	if toEmailsStr == "" {
		logger.Log.Fatal("No TO_EMAILS specified")
	}

	toEmailsList := strings.Split(toEmailsStr, ",")

	var toEmailsPtrs []*string

	for _, email := range toEmailsList {
		toEmailsPtrs = append(toEmailsPtrs, &email)
	}

	return toEmailsPtrs
}

func fetchSubject(update parser.CertificateUpdate) string {
	writer := new(bytes.Buffer)
	subjectTemplate, err := template.New("subject").Parse(subjectTemplate)

	if err != nil {
		logger.Log.Fatal("Error parsing subject template")
	}

	err = subjectTemplate.ExecuteTemplate(writer, "subject", map[string]interface{}{
		"Filename": update.Filename,
	})

	if err != nil {
		logger.Log.Fatal("Error executing subject template")
	}

	str := writer.String()
	str = strings.ReplaceAll(str, "\n", "")

	return str
}
