package parser

import (
	"testing"

	"ittconv/internal/timecode"
)

func TestParseITT(t *testing.T) {
	ittSource := `<?xml version="1.0" encoding="UTF-8"?>
<tt xmlns="http://www.w3.org/ns/ttml"
    xmlns:tts="http://www.w3.org/ns/ttml#styling"
    xmlns:ttp="http://www.w3.org/ns/ttml#parameter"
    xml:lang="en-US"
    ttp:timeBase="smpte"
    ttp:frameRate="24">
  <head>
    <styling>
      <style xml:id="s1" tts:color="red"/>
      <style xml:id="s2" tts:fontStyle="italic"/>
    </styling>
    <layout>
      <region xml:id="r1" tts:origin="10% 10%" tts:extent="80% 80%"/>
    </layout>
  </head>
  <body region="r1" style="s1">
    <div>
      <p begin="00:00:01:00" end="00:00:02:12" style="s2">Hello <span style="s1">World</span>!</p>
      <p begin="00:00:03:00" end="00:00:04:00">Line<br/>Break</p>
    </div>
  </body>
</tt>`

	doc, err := ParseITT(ittSource)
	if err != nil {
		t.Fatalf("ParseITT failed: %v", err)
	}

	if doc.Lang != "en-US" {
		t.Errorf("Expected Lang 'en-US', got %s", doc.Lang)
	}
	if doc.FrameRate != "24" {
		t.Errorf("Expected FrameRate '24', got %s", doc.FrameRate)
	}

	if len(doc.Styles) != 2 {
		t.Errorf("Expected 2 styles, got %d", len(doc.Styles))
	}
	if _, ok := doc.Styles["s1"]; !ok {
		t.Errorf("Expected style s1 to exist")
	}
	if _, ok := doc.Styles["s2"]; !ok {
		t.Errorf("Expected style s2 to exist")
	}

	if len(doc.Regions) != 1 {
		t.Errorf("Expected 1 region, got %d", len(doc.Regions))
	}
	if _, ok := doc.Regions["r1"]; !ok {
		t.Errorf("Expected region r1 to exist")
	}

	if len(doc.Cues) != 2 {
		t.Fatalf("Expected 2 cues, got %d", len(doc.Cues))
	}

	// Test first cue
	c1 := doc.Cues[0]
	fr24, _ := timecode.NewFrameRate("24")
	begin1, _ := timecode.ParseSMPTETimecode("00:00:01:00")
	beginRat1, _ := begin1.ToMilliseconds(fr24)

	if c1.Begin.Cmp(beginRat1) != 0 {
		t.Errorf("Cue 1: Expected Begin %s, got %s", beginRat1.String(), c1.Begin.String())
	}

	// SMPTE: 00:00:02:12 at 24fps -> (2 + 12/24) * 1000 = 2500ms
	end2Frames := timecode.SMPTETimecode{Hours: 0, Minutes: 0, Seconds: 2, Frames: 12}
	expectedEnd2, _ := end2Frames.ToMilliseconds(fr24)

	if c1.End.Cmp(expectedEnd2) != 0 {
		t.Errorf("Cue 1: Expected End %s, got %s", expectedEnd2.String(), c1.End.String())
	}
	if c1.RegionID != "r1" {
		t.Errorf("Cue 1: Expected RegionID 'r1', got %s", c1.RegionID)
	}
	if len(c1.StyleIDs) != 2 || c1.StyleIDs[0] != "s2" || c1.StyleIDs[1] != "s1" {
		t.Errorf("Cue 1: Expected StyleIDs ['s2', 's1'], got %v", c1.StyleIDs)
	}
	if c1.Content != "Hello World!" {
		t.Errorf("Cue 1: Expected Content 'Hello World!', got %s", c1.Content)
	}

	// Test second cue
	c2 := doc.Cues[1]
	begin2, _ := timecode.ParseSMPTETimecode("00:00:03:00")
	beginRat2, _ := begin2.ToMilliseconds(fr24)

	if c2.Begin.Cmp(beginRat2) != 0 {
		t.Errorf("Cue 2: Expected Begin %s, got %s", beginRat2.String(), c2.Begin.String())
	}
	// SMPTE: 00:00:04:00 at 24fps -> 4 * 1000 = 4000ms
	end2Frames_2 := timecode.SMPTETimecode{Hours: 0, Minutes: 0, Seconds: 4, Frames: 0}
	expectedEnd2_2, _ := end2Frames_2.ToMilliseconds(fr24)

	if c2.End.Cmp(expectedEnd2_2) != 0 {
		t.Errorf("Cue 2: Expected End %s, got %s", expectedEnd2_2.String(), c2.End.String())
	}
	if c2.Content != "Line\nBreak" {
		t.Errorf("Cue 2: Expected Content 'Line\nBreak', got %s", c2.Content)
	}
}
