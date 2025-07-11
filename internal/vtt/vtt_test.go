package vtt

import (
	"math/big"
	"testing"

	"github.com/mediafellows/ittconv/internal/parser"

	"github.com/google/go-cmp/cmp"
)

func TestToVTT(t *testing.T) {
	doc := &parser.ITTDocument{
		Cues: []parser.Cue{
			{
				ID:      "cue2",
				Begin:   big.NewRat(3000, 1), // 3000ms
				End:     big.NewRat(4000, 1), // 4000ms
				Content: "Line<br/>Break",
			},
			{
				ID:      "cue1",
				Begin:   big.NewRat(1000, 1), // 1000ms
				End:     big.NewRat(2500, 1), // 2500ms
				Content: "Hello World!",
			},
		},
	}

	expectedVTT := `WEBVTT

1
00:00:01.000 --> 00:00:02.500
Hello World!

2
00:00:03.000 --> 00:00:04.000
Line
Break
`

	vtt, err := ToVTT(doc)
	if err != nil {
		t.Fatalf("ToVTT failed: %v", err)
	}

	if diff := cmp.Diff(expectedVTT, vtt); diff != "" {
		t.Errorf("VTT output mismatch (-want +got):\n%s", diff)
	}
}
