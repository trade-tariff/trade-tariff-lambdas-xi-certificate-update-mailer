package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/fetcher"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/logger"
)

func setup() {
	logger.Log = &logger.MockLogger{}
}

func loadTestFile(filename string) *fetcher.XmlFile {
	content, _ := os.ReadFile(filename)

	return &fetcher.XmlFile{
		Content:  content,
		Key:      filename,
		LoadedOn: "2023-09-02",
	}
}

func TestParseXML_ValidCertificates(t *testing.T) {
	setup()
	file := loadTestFile("testdata/certificates.xml")

	result := ParseXML(file)

	filename := result.Filename
	loadedOn := result.LoadedOn
	change := result.CertificateChanges[0]
	certificate := change.Certificate
	period := change.CertificatePeriod
	description := change.CertificateDescription

	assert.Equal(t, "certificates.xml", filename)
	assert.Equal(t, "2023-09-02", loadedOn)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.CertificateChanges))
	assert.NotNil(t, change)
	assert.NotNil(t, certificate)
	assert.NotNil(t, period)
	assert.NotNil(t, description)
	assert.Equal(t, "created", certificate.UpdateType)
	assert.Equal(t, "X", certificate.TypeCode)
	assert.Equal(t, "808", certificate.Code)
	assert.Equal(t, "2023-08-05", certificate.ValidityStartDate)
	assert.Equal(t, "", certificate.ValidityEndDate)
	assert.Equal(t, "created", period.UpdateType)
	assert.Equal(t, "4666", period.SID)
	assert.Equal(t, "X", period.TypeCode)
	assert.Equal(t, "808", period.Code)
	assert.Equal(t, "2023-08-05", period.ValidityStartDate)
	assert.Equal(t, "", period.ValidityEndDate)
	assert.Equal(t, "created", description.UpdateType)
	assert.Equal(t, "4666", description.SID)
	assert.Equal(t, "EN", description.LangID)
	assert.Equal(t, "X", description.TypeCode)
	assert.Equal(t, "808", description.Code)
}

func TestParseXML_EmptyFile(t *testing.T) {
	setup()
	file := loadTestFile("testdata/empty.xml")

	result := ParseXML(file)

	assert.Empty(t, result.CertificateChanges)
}

func TestParseXML_NoCertificates(t *testing.T) {
	setup()
	file := loadTestFile("testdata/nocertificates.xml")

	result := ParseXML(file)

	assert.Empty(t, result.CertificateChanges)
}
