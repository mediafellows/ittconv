package vtt

import (
	"math/big"
	"testing"

	"github.com/mediafellows/ittconv/internal/parser"

	"github.com/google/go-cmp/cmp"
)

func TestToVTT(t *testing.T) {
	doc := &parser.ITTDocument{
		Styles: map[string]parser.Style{
			"style.default": {
				ID:    "style.default",
				Color: "white",
			},
			"style.D1AD2C": {
				ID:    "style.D1AD2C",
				Color: "#D1AD2C",
			},
			"style.D1AD2C.italic": {
				ID:        "style.D1AD2C.italic",
				Color:     "#D1AD2C",
				FontStyle: "italic",
			},
		},
		Regions: map[string]parser.Region{
			"bottom": {
				ID:           "bottom",
				TextAlign:    "center",
				DisplayAlign: "after",
			},
			"top": {
				ID:           "top",
				TextAlign:    "center",
				DisplayAlign: "before",
			},
		},
		Cues: []parser.Cue{
			{
				ID:       "cue1",
				Begin:    big.NewRat(1000, 1),
				End:      big.NewRat(2500, 1),
				RegionID: "bottom",
				Content:  `Hello &amp; World!`,
			},
			{
				ID:       "cue2",
				Begin:    big.NewRat(3000, 1),
				End:      big.NewRat(4000, 1),
				RegionID: "bottom",
				Content:  `Line<br/>Break`,
			},
			{
				ID:       "cue3",
				Begin:    big.NewRat(9001, 2),
				End:      big.NewRat(6000, 1),
				RegionID: "top",
				Content:  `<span style="style.D1AD2C">Plain </span><span style="style.D1AD2C.italic">italic</span>`,
			},
		},
	}

	expectedVTT := "\ufeffWEBVTT\r\n\r\n" +
		"00:00:01.000 --> 00:00:02.500 align:center line:-1\r\n" +
		"Hello &amp; World!\r\n\r\n" +
		"00:00:03.000 --> 00:00:04.000 align:center line:-2\r\n" +
		"Line\r\n" +
		"Break\r\n\r\n" +
		"00:00:04.501 --> 00:00:06.000 align:center line:0\r\n" +
		"<c.colorD1AD2C>Plain <i>italic</i></c>"

	vtt, err := ToVTT(doc)
	if err != nil {
		t.Fatalf("ToVTT failed: %v", err)
	}

	if diff := cmp.Diff(expectedVTT, vtt); diff != "" {
		t.Errorf("VTT output mismatch (-want +got):\n%s", diff)
	}
}
