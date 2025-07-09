package parser

import (
	"ittconv/internal/timecode"
	"math/big"
)

// ITTDocument represents the root of an iTunes Timed Text file.
type ITTDocument struct {
	Lang      string
	TimeBase  string
	FrameRate string
	Styles    map[string]Style
	Regions   map[string]Region
	Cues      []Cue
}

// Style represents a TTML style definition.
type Style struct {
	ID         string
	FontFamily string
	FontSize   string
	FontWeight string
	FontStyle  string
	Color      string
	// Add other styling attributes as needed
}

// Region represents a TTML region definition.
type Region struct {
	ID           string
	Origin       string
	Extent       string
	TextAlign    string
	DisplayAlign string
}

// Cue represents a single subtitle entry.
type Cue struct {
	ID            string
	Begin         *big.Rat
	End           *big.Rat
	BeginTimecode *timecode.SMPTETimecode // Temporary storage for SMPTE timecode
	EndTimecode   *timecode.SMPTETimecode // Temporary storage for SMPTE timecode
	RegionID      string
	StyleIDs      []string
	Content       string
}
