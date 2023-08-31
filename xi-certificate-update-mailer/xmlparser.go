package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
)

const (
	UpdateOperation  = "updated"
	DestroyOperation = "destroyed"
	CreateOperation  = "created"
)

func operationType(value string) string {
	switch value {
	case "1":
		return UpdateOperation
	case "2":
		return DestroyOperation
	case "3":
		return CreateOperation
	default:
		return "unknown"
	}
}

type CertificateUpdate struct {
	CertificateChanges []CertificateChange
}

func (c CertificateUpdate) Today() string {
	now := time.Now()

	return now.Format("2 Jan 2006")
}

type CertificateChange struct {
	Certificate            Certificate
	CertificatePeriod      CertificateDescriptionPeriod
	CertificateDescription CertificateDescription
}

func (c CertificateChange) FullCertificateCode() string {
	return c.Certificate.TypeCode + c.Certificate.Code
}

func (c CertificateChange) ValidityStartDate() string {
	if c.Certificate.ValidityStartDate != "" {
		return c.Certificate.ValidityStartDate
	} else if c.CertificatePeriod.ValidityStartDate != "" {
		return c.CertificatePeriod.ValidityStartDate
	} else {
		return ""
	}
}

func (c CertificateChange) ValidityEndDate() string {
	if c.Certificate.ValidityEndDate != "" {
		return c.Certificate.ValidityEndDate
	} else if c.CertificatePeriod.ValidityEndDate != "" {
		return c.CertificatePeriod.ValidityEndDate
	} else {
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
	if c.Certificate.UpdateType != "" {
		return c.Certificate.UpdateType
	} else if c.CertificatePeriod.UpdateType != "" {
		return c.CertificatePeriod.UpdateType
	} else if c.CertificateDescription.UpdateType != "" {
		return c.CertificateDescription.UpdateType
	} else {
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

type CertificateDescriptionPeriod struct {
	UpdateType        string
	SID               string
	TypeCode          string
	Code              string
	ValidityStartDate string
	ValidityEndDate   string
}

type CertificateDescription struct {
	UpdateType  string
	SID         string
	LangID      string
	TypeCode    string
	Code        string
	Description string
}

func getElementText(node *xmlquery.Node, elementName string) string {
	el := node.SelectElement(elementName)
	if el != nil {
		return el.InnerText()
	}
	return ""
}

func populateCertificate(node *xmlquery.Node) Certificate {
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

func populateCertificateDescriptionPeriod(node *xmlquery.Node) CertificateDescriptionPeriod {
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

func populateCertificateDescription(node *xmlquery.Node) CertificateDescription {
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

func FetchXML() []byte {
	xmlContent, err := os.ReadFile("test.xml")
	if err != nil {
		log.Fatalf("Failed reading XML: %v", err)
	}

	return xmlContent
}

func ParseXML(xml []byte) CertificateUpdate {
	var certificateChanges []CertificateChange

	doc, err := xmlquery.Parse(strings.NewReader(string(xml)))

	if err != nil {
		panic(err)
	}

	transactions := xmlquery.Find(doc, "//env:envelope/env:transaction")

	for _, t := range transactions {
		certificateNode := xmlquery.FindOne(t, ".//oub:certificate")
		descriptionPeriodNode := xmlquery.FindOne(t, ".//oub:certificate.description.period")
		descriptionNode := xmlquery.FindOne(t, ".//oub:certificate.description")

		if certificateNode != nil || descriptionPeriodNode != nil || descriptionNode != nil {
			certificateChanges = append(
				certificateChanges,
				CertificateChange{
					Certificate:            populateCertificate(certificateNode),
					CertificatePeriod:      populateCertificateDescriptionPeriod(descriptionPeriodNode),
					CertificateDescription: populateCertificateDescription(descriptionNode),
				},
			)
		}

	}

	return CertificateUpdate{CertificateChanges: certificateChanges}
}
