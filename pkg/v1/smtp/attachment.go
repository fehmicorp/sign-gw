package smtp

import (
	"log"

	"github.com/fehmicorp/sign-gw/pkg/v1/config"
)

// ProcessAttachments is called after parsing the MIME message.
func ProcessAttachments(email *config.Email) error {

	if len(email.Attachments) == 0 {
		return nil
	}

	log.Printf("Attachments Found : %d", len(email.Attachments))

	for _, a := range email.Attachments {

		log.Printf(
			"Attachment : %s (%s) %d bytes",
			a.FileName,
			a.ContentType,
			len(a.Data),
		)

		// Future:
		// • Virus Scan
		// • DLP
		// • Block Extensions
		// • Compress Images
		// • Save Metadata
	}

	return nil
}
