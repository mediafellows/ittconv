package ittconv

import (
	"github.com/mediafellows/ittconv/internal/parser"
	"github.com/mediafellows/ittconv/internal/ttml"
	"github.com/mediafellows/ittconv/internal/vtt"
)

// ToTTML converts an ITT source string to a TTML formatted string.
func ToTTML(ittSource string) (string, error) {
	doc, err := parser.ParseITT(ittSource)
	if err != nil {
		return "", err
	}
	return ttml.ToTTML(doc)
}

// ToVTT converts an ITT source string to a WebVTT formatted string.
func ToVTT(ittSource string) (string, error) {
	doc, err := parser.ParseITT(ittSource)
	if err != nil {
		return "", err
	}
	return vtt.ToVTT(doc)
}
