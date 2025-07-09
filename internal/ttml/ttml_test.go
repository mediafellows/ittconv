package ttml

import (
	"ittconv/internal/parser"
	"math/big"
	"strings"
	"testing"
)

func TestToTTML(t *testing.T) {
	doc := &parser.ITTDocument{
		Lang: "en-US",
		Styles: map[string]parser.Style{
			"s1": {ID: "s1", Color: "red"},
			"s2": {ID: "s2", FontStyle: "italic"},
		},
		Regions: map[string]parser.Region{
			"r1": {ID: "r1", Origin: "10% 10%"},
		},
		Cues: []parser.Cue{
			{
				Begin:    big.NewRat(1000, 1),
				End:      big.NewRat(2500, 1),
				Content:  "Hello World!",
				RegionID: "r1",
				StyleIDs: []string{"s1", "s2"},
			},
		},
	}

	expectedTTMLFragment := `<p begin="00:00:01.000" end="00:00:02.500" style="s1 s2" region="r1">Hello World!</p>`

	ttml, err := ToTTML(doc)
	if err != nil {
		t.Fatalf("ToTTML failed: %v", err)
	}

	if !strings.Contains(ttml, expectedTTMLFragment) {
		t.Errorf("Expected TTML to contain fragment:\n%s\n\nGot TTML:\n%s", expectedTTMLFragment, ttml)
	}

	// Also check for a style definition
	expectedStyleFragment := `<style xml:id="s1" tts:color="red"></style>`
	if !strings.Contains(ttml, expectedStyleFragment) {
		t.Errorf("Expected TTML to contain style fragment:\n%s\n\nGot TTML:\n%s", expectedStyleFragment, ttml)
	}
}
