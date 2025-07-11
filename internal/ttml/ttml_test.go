package ttml

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mediafellows/ittconv/internal/parser"
)

func TestToTTML(t *testing.T) {
	// First, parse the valid input fixture to get a document object.
	ittSource, err := ioutil.ReadFile("../../testdata/valid_input.itt")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}
	doc, err := parser.ParseITT(string(ittSource))
	if err != nil {
		t.Fatalf("Failed to parse ITT for TTML test: %v", err)
	}

	// Now, convert the parsed document to TTML.
	ttmlOutput, err := ToTTML(doc)
	if err != nil {
		t.Fatalf("ToTTML failed: %v", err)
	}

	// Check for expected TTML content.
	// 1. Check for correct time format conversion (first cue).
	//    00:00:03:12 @ 24fps = 3 + 12/24 = 3.5s = 3.500
	if !strings.Contains(ttmlOutput, `begin="00:00:01.000" end="00:00:03.500"`) {
		t.Errorf("Expected correct time conversion. Looked for `begin=\"00:00:01.000\" end=\"00:00:03.500\"`.\nGot:\n%s", ttmlOutput)
	}

	// 2. Check for preserved content with inline styles.
	if !strings.Contains(ttmlOutput, `<span style="s1">second</span>`) {
		t.Errorf("Expected inline span with style. Looked for `<span style=\"s1\">second</span>`.\nGot:\n%s", ttmlOutput)
	}

	// 3. Check for preserved line breaks.
	if !strings.Contains(ttmlOutput, "A third one<br/>with a line break.") {
		t.Errorf("Expected preserved line break. Looked for `A third one<br/>with a line break.`.\nGot:\n%s", ttmlOutput)
	}

	// 4. Check for correct XML header and root element.
	if !strings.HasPrefix(ttmlOutput, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML header prefix.")
	}
	if !strings.Contains(ttmlOutput, `<tt xmlns="http://www.w3.org/ns/ttml"`) {
		t.Error("Expected tt root element with ttml namespace.")
	}
}
