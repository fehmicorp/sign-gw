package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
	"github.com/fehmicorp/sign-gw/pkg/v1/smtp"
)

func test() {
	emlDir := filepath.Join(".", "data", "eml")
	inputPath := filepath.Join(emlDir, "test.eml")
	outputPath := filepath.Join(emlDir, "test_processed.eml")

	log.Printf("Reading EML file from: %s", inputPath)

	rawBytes, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Failed to read input EML: %v", err)
	}

	parsedEmail, err := smtp.ParseMessage(rawBytes)
	if err != nil {
		log.Fatalf("Failed to parse message: %v", err)
	}

	// Fetch dynamic HTML signature from LDAP & Template
	htmlSignature, err := smtp.HTMLSignature(parsedEmail)
	fmt.Printf("Html Signature: %s\n", htmlSignature)
	textSignature := smtp.HTMLToText(htmlSignature)
	fmt.Printf("Text Signature: %s\n", textSignature)
	// Construct email config object for Build
	emailConfig := &config.Email{
		Raw:          rawBytes,
		EnvelopeFrom: parsedEmail.EnvelopeFrom,
		EnvelopeTo:   parsedEmail.EnvelopeTo,
		Subject:      parsedEmail.Subject,
		HTML:         htmlSignature,
		Text:         textSignature,
	}

	// Process message and perform replacement of %%SIGN%%
	log.Println("Searching for %%SIGN%% and injecting signature...")
	processedBytes, err := smtp.Build(emailConfig)
	if err != nil {
		log.Fatalf("Signature injection failed: %v", err)
	}

	// Save the processed result
	err = os.WriteFile(outputPath, processedBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to save processed EML file: %v", err)
	}

	log.Printf("Successfully injected signature and saved output to: %s", outputPath)
}
