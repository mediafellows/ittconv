package ttml

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"ittconv/internal/parser"

	"github.com/google/go-cmp/cmp"
)

// TestSubtitlePairs iterates through every *.itt / *.ttml pair found in
// testdata/pairs/ and ensures the converter produces exactly the expected
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
		wantPath := filepath.Join(pairsDir, base+".ttml")

		ittData, err := ioutil.ReadFile(ittPath)
		if err != nil {
			t.Fatalf("failed to read %s: %v", ittPath, err)
		}

		wantData, err := ioutil.ReadFile(wantPath)
		if err != nil {
			t.Fatalf("missing matching .ttml file for %s: %v", ittPath, err)
		}

		doc, err := parser.ParseITT(string(ittData))
		if err != nil {
			t.Fatalf("failed to parse %s: %v", ittPath, err)
		}

		got, err := ToTTML(doc)
		if err != nil {
			t.Fatalf("failed to convert %s: %v", ittPath, err)
		}

		// Direct comparison – indentation and attribute ordering must match.
		if diff := cmp.Diff(string(wantData), got); diff != "" {
			t.Errorf("conversion mismatch for %s (-want +got):\n%s", filepath.Base(ittPath), diff)
		}
	}
}
