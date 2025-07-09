package vtt

import (
	"fmt"
	"ittconv/internal/parser"
	"math/big"
	"sort"
	"strings"
)

// ToVTT converts an ITTDocument to a VTT formatted string.
func ToVTT(doc *parser.ITTDocument) (string, error) {
	var sb strings.Builder

	sb.WriteString("WEBVTT\n\n")

	// Sort cues by begin time
	sort.Slice(doc.Cues, func(i, j int) bool {
		return doc.Cues[i].Begin.Cmp(doc.Cues[j].Begin) < 0
	})

	for _, cue := range doc.Cues {
		beginStr, err := formatVTTTimestamp(cue.Begin)
		if err != nil {
			return "", fmt.Errorf("error formatting begin timestamp for cue %s: %w", cue.ID, err)
		}
		endStr, err := formatVTTTimestamp(cue.End)
		if err != nil {
			return "", fmt.Errorf("error formatting end timestamp for cue %s: %w", cue.ID, err)
		}

		sb.WriteString(fmt.Sprintf("%s --> %s\n", beginStr, endStr))
		sb.WriteString(cue.Content)
		sb.WriteString("\n\n")
	}

	return sb.String(), nil
}

// formatVTTTimestamp converts a big.Rat (in milliseconds) to a VTT timestamp string (HH:MM:SS.mmm).
func formatVTTTimestamp(ms *big.Rat) (string, error) {
	if ms == nil {
		return "00:00:00.000", nil
	}

	// To convert a rational number a/b to an integer, we can compute (a * 1000) / b
	// and then take the integer part. However, we are dealing with milliseconds, so we
	// can just compute a/b.
	msInt := new(big.Int).Quo(ms.Num(), ms.Denom()).Int64()

	hours := msInt / 3600000
	msInt %= 3600000
	minutes := msInt / 60000
	msInt %= 60000
	seconds := msInt / 1000
	milliseconds := msInt % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds), nil
}
