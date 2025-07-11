package vtt

import (
	"bytes"
	"sort"
	"strings"

	"github.com/mediafellows/ittconv/internal/parser"
	"github.com/mediafellows/ittconv/internal/ttml"

	"github.com/asticode/go-astisub"
)

// ToVTT converts an ITTDocument to a VTT formatted string.
// It does so by first converting the document to a temporary TTML string,
// and then uses the go-astisub library to perform a high-fidelity
// conversion from TTML to VTT, preserving styling and region information.
func ToVTT(doc *parser.ITTDocument) (string, error) {
	// Step 1: Convert our internal ITTDocument to a TTML string.
	ttmlString, err := ttml.ToTTML(doc)
	if err != nil {
		return "", err
	}

	// Step 2: Use astisub to read the TTML from a string reader.
	subs, err := astisub.ReadFromTTML(strings.NewReader(ttmlString))
	if err != nil {
		return "", err
	}

	// Sort cues by timestamp to ensure deterministic output.
	sort.Slice(subs.Items, func(i, j int) bool {
		return subs.Items[i].StartAt < subs.Items[j].StartAt
	})

	// Step 3: Write the subtitles to a WebVTT format in a buffer.
	var buf bytes.Buffer
	if err := subs.WriteToWebVTT(&buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
