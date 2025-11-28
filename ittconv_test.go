package ittconv

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/asticode/go-astisub"
)

func TestToVTT(t *testing.T) {
	// Read the valid ITT fixture
	ittSource, err := ioutil.ReadFile("testdata/valid_input.itt")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	// Convert to VTT
	vttOutput, err := ToVTT(string(ittSource))
	if err != nil {
		t.Fatalf("ToVTT failed: %v", err)
	}

	// Basic checks for VTT content
	if !strings.HasPrefix(vttOutput, "WEBVTT") {
		t.Error("Expected VTT output to start with WEBVTT")
	}
	if !strings.Contains(vttOutput, "-->") {
		t.Error("Expected VTT output to contain time cues")
	}
	if !strings.Contains(vttOutput, "A third one") {
		t.Error("Expected VTT output to contain subtitle text")
	}
}

func TestToTTML(t *testing.T) {
	// Read the valid ITT fixture
	ittSource, err := ioutil.ReadFile("testdata/valid_input.itt")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	// Convert to TTML
	ttmlOutput, err := ToTTML(string(ittSource))
	if err != nil {
		t.Fatalf("ToTTML failed: %v", err)
	}

	// Basic checks for TTML content
	if !strings.HasPrefix(ttmlOutput, "<?xml version") {
		t.Error("Expected TTML output to start with XML declaration")
	}
	if !strings.Contains(ttmlOutput, "<tt xmlns") {
		t.Error("Expected TTML output to contain the tt root element")
	}
	if !strings.Contains(ttmlOutput, "<p begin=") {
		t.Error("Expected TTML output to contain p elements with begin attribute")
	}
}

func TestConversionChainFixtures(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		wantTTML []string
		wantVTT  []string
	}{
		{
			name: "AngleBracketEscapes",
			path: "testdata/test_brackets.itt",
			wantTTML: []string{
				"&lt;FLASH SALE&gt;",
				"&lt;promo&gt;",
			},
			wantVTT: []string{
				"&lt;FLASH SALE>",
				"&lt;promo>",
			},
		},
		{
			name: "AmpersandEscapes",
			path: "testdata/test_ampersands.itt",
			wantTTML: []string{
				"AT&amp;T",
				"Butter &amp; Jam",
			},
			wantVTT: []string{
				"AT&amp;T",
				"Butter &amp; Jam",
			},
		},
		{
			name: "QuotesAndEntities",
			path: "testdata/test_quotes.itt",
			wantTTML: []string{
				"&#34;run&#34;",
				"&amp; vanished",
				"© and —",
			},
			wantVTT: []string{
				`"run"`,
				"&amp; vanished",
				"©",
				"—",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ittSource, err := ioutil.ReadFile(tc.path)
			if err != nil {
				t.Fatalf("Failed to read test fixture %s: %v", tc.path, err)
			}

			ttmlOutput, err := ToTTML(string(ittSource))
			if err != nil {
				t.Fatalf("ToTTML failed for %s: %v", tc.path, err)
			}
			if len(ttmlOutput) == 0 {
				t.Fatalf("TTML output is empty for %s", tc.path)
			}

			for _, snippet := range tc.wantTTML {
				if !strings.Contains(ttmlOutput, snippet) {
					t.Fatalf("TTML output for %s missing snippet %q", tc.path, snippet)
				}
			}

			subs, err := astisub.ReadFromTTML(strings.NewReader(ttmlOutput))
			if err != nil {
				t.Fatalf("go-astisub failed to parse TTML output for %s: %v", tc.path, err)
			}

			var buf bytes.Buffer
			if err := subs.WriteToWebVTT(&buf); err != nil {
				t.Fatalf("go-astisub failed to convert TTML to WebVTT for %s: %v", tc.path, err)
			}

			vttOutput := buf.String()
			if !strings.HasPrefix(vttOutput, "WEBVTT") {
				t.Fatalf("Converted WebVTT output for %s does not start with WEBVTT header", tc.path)
			}
			if !strings.Contains(vttOutput, "-->") {
				t.Fatalf("Converted WebVTT output for %s missing cue timestamps", tc.path)
			}
			for _, snippet := range tc.wantVTT {
				if !strings.Contains(vttOutput, snippet) {
					t.Fatalf("WebVTT output for %s missing snippet %q", tc.path, snippet)
				}
			}
		})
	}
}
