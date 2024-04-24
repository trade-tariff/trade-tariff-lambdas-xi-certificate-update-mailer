package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/fetcher"
	"github.com/trade-tariff/trade-tariff-lambdas-xi-certificate-update-mailer/logger"
)

const (
	updateOperation  = "updated"
	destroyOperation = "destroyed"
	createOperation  = "created"
)

type CertificateUpdate struct {
	CertificateChanges []CertificateChange
	Filename           string
	LoadedOn           string
}

type CertificateChange struct {
	Certificate            Certificate
	CertificatePeriod      CertificateDescriptionPeriod
	CertificateDescription CertificateDescription
}

func newCertificateChange(transaction *xmlquery.Node) *CertificateChange {
	certificateNode := xmlquery.FindOne(transaction, ".//oub:record/oub:certificate")
	descriptionPeriodNode := xmlquery.FindOne(transaction, ".//oub:record/oub:certificate.description.period")
	descriptionNode := xmlquery.FindOne(transaction, ".//oub:record/oub:certificate.description")
	hasCertificate := certificateNode != nil
	hasDescriptionPeriod := descriptionPeriodNode != nil
	hasDescription := descriptionNode != nil

	debug := os.Getenv("DEBUG") == "true"

	if debug {
		transactionMessageCount := len(xmlquery.Find(transaction, "./env:app.message"))
		allMessageKeys := map[string]struct{}{}

		for _, message := range xmlquery.Find(transaction, "./env:app.message") {
			for _, element := range xmlquery.Find(message, ".//*") {
				allMessageKeys[element.Data] = struct{}{}
			}
		}

		logger.Log.Info(
			"Found certificate change",
			logger.String("transactionMessageCount", fmt.Sprintf("%d", transactionMessageCount)),
			logger.String("hasCertificate", fmt.Sprintf("%t", hasCertificate)),
			logger.String("hasDescriptionPeriod", fmt.Sprintf("%t", hasDescriptionPeriod)),
			logger.String("hasDescription", fmt.Sprintf("%t", hasDescription)),
			logger.String("allMessageKeys", fmt.Sprintf("%v", allMessageKeys)),
		)
	}

	if hasCertificate || hasDescriptionPeriod || hasDescription {
		return &CertificateChange{
			Certificate:            newCertificate(certificateNode),
			CertificatePeriod:      newCertificateDescriptionPeriod(descriptionPeriodNode),
			CertificateDescription: newCertificateDescription(descriptionNode),
		}
	} else {
		return nil
	}
}

func (c CertificateChange) FullCertificateCode() string {
	return c.CertificateTypeCode() + c.CertificateCode()
}

func (c CertificateChange) ValidityStartDate() string {
	switch {
	case c.Certificate.ValidityStartDate != "":
		return c.Certificate.ValidityStartDate
	case c.CertificatePeriod.ValidityStartDate != "":
		return c.CertificatePeriod.ValidityStartDate
	default:
		return ""
	}
}

func (c CertificateChange) ValidityEndDate() string {
	switch {
	case c.Certificate.ValidityEndDate != "":
		return c.Certificate.ValidityEndDate
	case c.CertificatePeriod.ValidityEndDate != "":
		return c.CertificatePeriod.ValidityEndDate
	default:
		return ""
	}
}

func (c CertificateChange) Description() string {
	if c.CertificateDescription.Description != "" {
		return c.CertificateDescription.Description
	} else {
		return ""
	}
}

func (c CertificateChange) Operation() string {
	switch {
	case c.Certificate.UpdateType != "":
		return c.Certificate.UpdateType
	case c.CertificatePeriod.UpdateType != "":
		return c.CertificatePeriod.UpdateType
	case c.CertificateDescription.UpdateType != "":
		return c.CertificateDescription.UpdateType
	default:
		return "unknown"
	}
}

func (c CertificateChange) CertificateTypeCode() string {
	if c.Certificate.TypeCode != "" {
		return c.Certificate.TypeCode
	} else if c.CertificatePeriod.TypeCode != "" {
		return c.CertificatePeriod.TypeCode
	} else if c.CertificateDescription.TypeCode != "" {
		return c.CertificateDescription.TypeCode
	} else {
		return ""
	}
}

func (c CertificateChange) CertificateCode() string {
	if c.Certificate.Code != "" {
		return c.Certificate.Code
	} else if c.CertificatePeriod.Code != "" {
		return c.CertificatePeriod.Code
	} else if c.CertificateDescription.Code != "" {
		return c.CertificateDescription.Code
	} else {
		return ""
	}
}

type Certificate struct {
	UpdateType        string
	TypeCode          string
	Code              string
	ValidityStartDate string
	ValidityEndDate   string
}

func newCertificate(node *xmlquery.Node) Certificate {
	if node == nil {
		return Certificate{}
	}

	return Certificate{
		UpdateType:        operationType(getElementText(node.Parent, "oub:update.type")),
		TypeCode:          getElementText(node, "oub:certificate.type.code"),
		Code:              getElementText(node, "oub:certificate.code"),
		ValidityStartDate: getElementText(node, "oub:validity.start.date"),
		ValidityEndDate:   getElementText(node, "oub:validity.end.date"),
	}
}

type CertificateDescriptionPeriod struct {
	UpdateType        string
	SID               string
	TypeCode          string
	Code              string
	ValidityStartDate string
	ValidityEndDate   string
}

func newCertificateDescriptionPeriod(node *xmlquery.Node) CertificateDescriptionPeriod {
	if node == nil {
		return CertificateDescriptionPeriod{}
	}

	return CertificateDescriptionPeriod{
		UpdateType:        operationType(getElementText(node.Parent, "oub:update.type")),
		SID:               getElementText(node, "oub:certificate.description.period.sid"),
		TypeCode:          getElementText(node, "oub:certificate.type.code"),
		Code:              getElementText(node, "oub:certificate.code"),
		ValidityStartDate: getElementText(node, "oub:validity.start.date"),
		ValidityEndDate:   getElementText(node, "oub:validity.end.date"),
	}
}

type CertificateDescription struct {
	UpdateType  string
	SID         string
	LangID      string
	TypeCode    string
	Code        string
	Description string
}

func newCertificateDescription(node *xmlquery.Node) CertificateDescription {
	if node == nil {
		return CertificateDescription{}
	}

	return CertificateDescription{
		UpdateType:  operationType(getElementText(node.Parent, "oub:update.type")),
		SID:         getElementText(node, "oub:certificate.description.period.sid"),
		LangID:      getElementText(node, "oub:language.id"),
		TypeCode:    getElementText(node, "oub:certificate.type.code"),
		Code:        getElementText(node, "oub:certificate.code"),
		Description: getElementText(node, "oub:description"),
	}
}

func ParseXML(file *fetcher.XmlFile) CertificateUpdate {
	var certificateChanges []CertificateChange

	doc, err := xmlquery.Parse(strings.NewReader(string(file.Content)))

	if err != nil {
		logger.Log.Fatal(
			"Error occurred while parsing xml",
			logger.String("filename", file.Filename()),
			logger.String("error", err.Error()),
		)
	}

	transactions := xmlquery.Find(doc, "//env:envelope/env:transaction")

	for _, t := range transactions {
		certificateChange := newCertificateChange(t)

		if certificateChange != nil {
			certificateChanges = append(certificateChanges, *certificateChange)
		}
	}

	return CertificateUpdate{
		CertificateChanges: certificateChanges,
		Filename:           file.Filename(),
		LoadedOn:           file.LoadedOn,
	}
}

func getElementText(node *xmlquery.Node, elementName string) string {
	el := node.SelectElement(elementName)
	if el != nil {
		return el.InnerText()
	}
	return ""
}

func operationType(value string) string {
	switch value {
	case "1":
		return updateOperation
	case "2":
		return destroyOperation
	case "3":
		return createOperation
	default:
		return "unknown"
	}
}
