package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	xml := FetchXML()
	certificateUpdate := ParseXML(xml)

	if len(certificateUpdate.CertificateChanges) > 0 {
		emailConfiguration := fetchEmailConfiguration(certificateUpdate)

		SendEmailSES(emailConfiguration)
	} else {
		log.Println("No certificate changes")
	}
}
