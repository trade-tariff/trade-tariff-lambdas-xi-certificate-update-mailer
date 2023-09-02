package emailer

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/stretchr/testify/assert"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/logger"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/parser"
)

type mockSESClient struct {
	input *ses.SendEmailInput
}

func (m *mockSESClient) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	m.input = input
	return &ses.SendEmailOutput{}, nil
}

func setup() {
	logger.Initialize()

	os.Setenv("FROM_EMAIL", "Online Trade Tariff Support <trade-tariff-support@enginegroup.com>")
	os.Setenv("TO_EMAILS", "Bob Dylan <recipient1@example.com>,Nina Simmone <recipient2@example.com>")
}

var update = parser.CertificateUpdate{
	Filename: "data/taric/certificates.xml",
	LoadedOn: "2023-09-02",
	CertificateChanges: []parser.CertificateChange{
		{
			Certificate: parser.Certificate{
				UpdateType:        "created",
				TypeCode:          "X",
				Code:              "808",
				ValidityStartDate: "2023-08-05",
				ValidityEndDate:   "",
			},
			CertificatePeriod: parser.CertificateDescriptionPeriod{
				UpdateType:        "created",
				SID:               "4666",
				TypeCode:          "X",
				Code:              "808",
				ValidityStartDate: "2023-08-05",
				ValidityEndDate:   "",
			},
			CertificateDescription: parser.CertificateDescription{
				UpdateType: "created",
				SID:        "4666",
				LangID:     "EN",
				TypeCode:   "X",
				Code:       "808",
			},
		},
	},
}

var expectedBody = `<p>Good morning,</p>

<p>
  The Northern Ireland tariff updates loaded on <b>2023-09-02</b> contain one or more certificate changes.
</p>


<p>Certificate <b>X808</b> has been created.</p>
<table cellspacing="0" cellpadding="0" style="width:80%">
  <tr>
    <td style="min-width:20%;vertical-align:top">Certificate code:</td>
    <td style="width:80%;vertical-align:top">X808</td>
  </tr>
  <tr>
    <td style="vertical-align:top">Valid from:</td>
    <td style="vertical-align:top">2023-08-05</td>
  </tr>
  <tr>
    <td style="vertical-align:top">Valid to:</td>
    <td style="vertical-align:top"></td>
  </tr>
  <tr>
    <td style="vertical-align:top">Description:</td>
    <td style="vertical-align:top"></td>
  </tr>
</table>

<p>
  To find out more, visit
  https://www.trade-tariff.service.gov.uk/xi/certificate_search?type=X&code=808
</p>


<hr>

`

func TestSend(t *testing.T) {
	setup()

	sesClient := &mockSESClient{}
	emailer := NewEmailer(sesClient)

	emailConfiguration := NewEmailConfiguration(update)
	emailer.Send(emailConfiguration)

	assert.Equal(t, emailConfiguration.fromEmail, *sesClient.input.Source)
	assert.Equal(t, emailConfiguration.subject, *sesClient.input.Message.Subject.Data)
	assert.Equal(t, emailConfiguration.body, *sesClient.input.Message.Body.Html.Data)
}

func TestNewEmailer(t *testing.T) {
	setup()

	config := NewEmailConfiguration(update)

	assert.Equal(t, "Online Trade Tariff Support <trade-tariff-support@enginegroup.com>", config.fromEmail)
	assert.Equal(t, "Certificate updates received in EU / XI data file data/taric/certificates.xml", config.subject)
	assert.Equal(t, expectedBody, config.body)
	assert.Equal(t, 2, len(config.toEmails))
	assert.Equal(t, "UTF-8", config.charset)
}
