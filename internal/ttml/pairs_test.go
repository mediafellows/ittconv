package ttml_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mediafellows/ittconv/internal/parser"
	"github.com/mediafellows/ittconv/internal/ttml"
	"github.com/mediafellows/ittconv/internal/vtt"

	"github.com/google/go-cmp/cmp"
)

// TestSubtitlePairs iterates through every *.itt pair found in testdata/pairs/
// and ensures the converter produces matching TTML *and* VTT outputs.
// output. When a mismatch occurs, a unified diff is printed to help narrow
// down the problem area.
func TestSubtitlePairs(t *testing.T) {
	// The pairs directory is two levels up from this package directory.
	pairsDir := filepath.Join("..", "..", "testdata", "pairs")

	ittFiles, err := filepath.Glob(filepath.Join(pairsDir, "*.itt"))
	if err != nil {
		t.Fatalf("failed to list .itt files in %s: %v", pairsDir, err)
	}
	if len(ittFiles) == 0 {
		t.Fatalf("no .itt files found in %s", pairsDir)
	}

	for _, ittPath := range ittFiles {
		base := strings.TrimSuffix(filepath.Base(ittPath), filepath.Ext(ittPath))
		wantTTML := filepath.Join(pairsDir, base+".ttml")
		wantVTT := filepath.Join(pairsDir, base+".vtt")

		ittData, err := ioutil.ReadFile(ittPath)
		if err != nil {
			t.Fatalf("failed to read %s: %v", ittPath, err)
		}

		wantTTMLData, err := ioutil.ReadFile(wantTTML)
		if err != nil {
			t.Fatalf("missing matching .ttml file for %s: %v", ittPath, err)
		}

		wantVTTData, err := ioutil.ReadFile(wantVTT)
		if err != nil {
			t.Fatalf("missing matching .vtt file for %s: %v", ittPath, err)
		}

		doc, err := parser.ParseITT(string(ittData))
		if err != nil {
			t.Fatalf("failed to parse %s: %v", ittPath, err)
		}

		gotTTML, err := ttml.ToTTML(doc)
		if err != nil {
			t.Fatalf("failed to convert %s: %v", ittPath, err)
		}

		// Direct comparison – indentation and attribute ordering must match.
		if diff := cmp.Diff(string(wantTTMLData), gotTTML); diff != "" {
			t.Errorf("TTML mismatch for %s (-want +got):\n%s", filepath.Base(ittPath), diff)
		}

		// Now VTT
		gotVTT, err := vtt.ToVTT(doc)
		if err != nil {
			t.Fatalf("failed to convert %s to VTT: %v", ittPath, err)
		}
		if diff := cmp.Diff(string(wantVTTData), gotVTT); diff != "" {
			t.Errorf("VTT mismatch for %s (-want +got):\n%s", filepath.Base(ittPath), diff)
		}
	}
}
