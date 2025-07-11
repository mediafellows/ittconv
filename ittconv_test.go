package ittconv

import (
	"io/ioutil"
	"strings"
	"testing"
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
