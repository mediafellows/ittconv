package parser

import (
	"io/ioutil"
	"math/big"
	"strings"
	"testing"

	"github.com/mediafellows/ittconv/internal/timecode"
)

func TestParseITT_Valid(t *testing.T) {
	ittSource, err := ioutil.ReadFile("../../testdata/valid_input.itt")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	doc, err := ParseITT(string(ittSource))
	if err != nil {
		t.Fatalf("ParseITT failed: %v", err)
	}

	if doc.Lang != "en-US" {
		t.Errorf("Expected Lang 'en-US', got '%s'", doc.Lang)
	}
	if doc.FrameRate != "24" {
		t.Errorf("Expected FrameRate '24', got '%s'", doc.FrameRate)
	}
	if len(doc.Styles) != 2 {
		t.Errorf("Expected 2 styles, got %d", len(doc.Styles))
	}
	if len(doc.Regions) != 1 {
		t.Errorf("Expected 1 region, got %d", len(doc.Regions))
	}
	if len(doc.Cues) != 3 {
		t.Fatalf("Expected 3 cues, got %d", len(doc.Cues))
	}

	// Spot check the second cue
	cue := doc.Cues[1]
	if !strings.Contains(cue.Content, "second") {
		t.Errorf("Expected second cue content to be correct, got '%s'", cue.Content)
	}
	if len(cue.StyleIDs) != 1 {
		t.Fatalf("Expected 1 style ID for the cue, got %d", len(cue.StyleIDs))
	}
	if cue.StyleIDs[0] != "s2" {
		t.Errorf("Expected style ID 's2' for the cue, got '%s'", cue.StyleIDs[0])
	}

	fr, _ := timecode.NewFrameRate("24")
	expectedBegin, _ := timecode.ParseSMPTETimecode("00:00:04:00")
	expectedBeginMs, _ := expectedBegin.ToMilliseconds(fr)
	if cue.Begin.Cmp(expectedBeginMs) != 0 {
		t.Errorf("Expected begin time %v, got %v", expectedBeginMs, cue.Begin)
	}
}

func TestParseITT_NoFrameRate(t *testing.T) {
	ittSource, err := ioutil.ReadFile("../../testdata/no_framerate.itt")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	_, err = ParseITT(string(ittSource))
	if err == nil {
		t.Fatal("Expected an error when parsing ITT with no frame rate, but got nil")
	}
	if !strings.Contains(err.Error(), "frameRate attribute missing") {
		t.Errorf("Expected error message about missing frame rate, but got: %v", err)
	}
}

func TestParseITT_InvalidTimeRange(t *testing.T) {
	ittSource, err := ioutil.ReadFile("../../testdata/invalid_time_range.itt")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	_, err = ParseITT(string(ittSource))
	if err == nil {
		t.Fatal("Expected an error when parsing ITT with invalid time range, but got nil")
	}
	if !strings.Contains(err.Error(), "is not less than end time") {
		t.Errorf("Expected error message about invalid time range, but got: %v", err)
	}
}

func TestParseITT_DivBeginOffset(t *testing.T) {
	ittSource, err := ioutil.ReadFile("../../testdata/shifted.itt")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	doc, err := ParseITT(string(ittSource))
	if err != nil {
		t.Fatalf("ParseITT failed: %v", err)
	}

	if doc.FrameRateMultiplierNum != 999 || doc.FrameRateMultiplierDen != 1000 {
		t.Fatalf("Expected frameRateMultiplier 999/1000, got %d/%d", doc.FrameRateMultiplierNum, doc.FrameRateMultiplierDen)
	}
	if len(doc.Cues) != 5 {
		t.Fatalf("Expected 5 cues, got %d", len(doc.Cues))
	}

	fr, _ := timecode.NewFrameRate("24")
	fr.Rat.Mul(fr.Rat, big.NewRat(999, 1000))

	shiftTc, _ := timecode.ParseSMPTETimecode("-00:59:59:00")
	shiftMs, _ := shiftTc.ToMilliseconds(fr)

	firstBeginTC, _ := timecode.ParseSMPTETimecode("01:05:27:21")
	firstBeginMs, _ := firstBeginTC.ToMilliseconds(fr)
	firstBeginMs.Add(firstBeginMs, shiftMs)

	firstCue := doc.Cues[0]
	if firstCue.Begin.Cmp(firstBeginMs) != 0 {
		t.Fatalf("Expected first cue begin %s, got %s", firstBeginMs.String(), firstCue.Begin.String())
	}

	lastBeginTC, _ := timecode.ParseSMPTETimecode("01:09:33:22")
	lastBeginMs, _ := lastBeginTC.ToMilliseconds(fr)
	lastBeginMs.Add(lastBeginMs, shiftMs)

	lastCue := doc.Cues[len(doc.Cues)-1]
	if lastCue.Begin.Cmp(lastBeginMs) != 0 {
		t.Fatalf("Expected last cue begin %s, got %s", lastBeginMs.String(), lastCue.Begin.String())
	}
}
