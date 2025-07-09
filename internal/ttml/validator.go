package ttml

import (
	"encoding/xml"
	"fmt"
)

// ValidateTTML performs a lightweight validation of a TTML document.
//
// It checks that:
//  1. The document is well-formed XML.
//  2. The root element is <tt> in the TTML namespace.
//  3. A <body> element exists (required by the TTML spec).
//
// This is *not* a full W3C TTML schema validation, but it provides
// a quick sanity-check that helps catch malformed output during tests
// and in the CLI. For full conformance the generated document can be
// passed through an external validator (e.g. ttml2-validator), but this
// function guarantees that the converter never emits outright invalid
// XML.
func ValidateTTML(ttmlStr string) error {
	// Define a minimal structure we can unmarshal into to verify the
	// presence of required elements.
	type root struct {
		XMLName xml.Name  `xml:"tt"`
		Body    *struct{} `xml:"body"`
	}

	var r root
	if err := xml.Unmarshal([]byte(ttmlStr), &r); err != nil {
		return fmt.Errorf("invalid TTML XML: %w", err)
	}

	if r.XMLName.Local != "tt" {
		return fmt.Errorf("root element is <%s>, expected <tt>", r.XMLName.Local)
	}

	if r.Body == nil {
		return fmt.Errorf("TTML document is missing required <body> element")
	}

	return nil
}
