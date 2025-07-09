package ttml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"ittconv/internal/parser"
	"math/big"
	"strings"
)

// ToTTML converts an ITTDocument to a standard TTML formatted string.
func ToTTML(doc *parser.ITTDocument) (string, error) {
	// Create a new structure that can be easily marshaled to XML
	type ttP struct {
		XMLName xml.Name `xml:"p"`
		Begin   string   `xml:"begin,attr"`
		End     string   `xml:"end,attr"`
		Content string   `xml:",innerxml"`
		Style   string   `xml:"style,attr,omitempty"`
		Region  string   `xml:"region,attr,omitempty"`
	}

	type ttDiv struct {
		XMLName xml.Name `xml:"div"`
		Ps      []ttP    `xml:"p"`
	}

	type ttBody struct {
		XMLName xml.Name `xml:"body"`
		Divs    []ttDiv  `xml:"div"`
		Region  string   `xml:"region,attr,omitempty"`
	}

	type ttStyle struct {
		XMLName    xml.Name `xml:"style"`
		ID         string   `xml:"xml:id,attr"`
		Color      string   `xml:"tts:color,attr,omitempty"`
		FontFamily string   `xml:"tts:fontFamily,attr,omitempty"`
		FontSize   string   `xml:"tts:fontSize,attr,omitempty"`
		FontStyle  string   `xml:"tts:fontStyle,attr,omitempty"`
		FontWeight string   `xml:"tts:fontWeight,attr,omitempty"`
	}

	type ttStyling struct {
		XMLName xml.Name  `xml:"styling"`
		Styles  []ttStyle `xml:"style"`
	}

	type ttRegion struct {
		XMLName      xml.Name `xml:"region"`
		ID           string   `xml:"xml:id,attr"`
		Origin       string   `xml:"origin,attr,omitempty"`
		Extent       string   `xml:"extent,attr,omitempty"`
		DisplayAlign string   `xml:"displayAlign,attr,omitempty"`
		TextAlign    string   `xml:"textAlign,attr,omitempty"`
	}

	type ttLayout struct {
		XMLName xml.Name   `xml:"layout"`
		Regions []ttRegion `xml:"region"`
	}

	type ttHead struct {
		XMLName xml.Name  `xml:"head"`
		Styling ttStyling `xml:"styling"`
		Layout  ttLayout  `xml:"layout"`
	}

	type ttRoot struct {
		XMLName  xml.Name `xml:"tt"`
		Xmlns    string   `xml:"xmlns,attr"`
		XmlnsTTP string   `xml:"xmlns:ttp,attr"`
		XmlnsTTS string   `xml:"xmlns:tts,attr"`
		TimeBase string   `xml:"ttp:timeBase,attr"`
		Lang     string   `xml:"xml:lang,attr"`
		Head     ttHead   `xml:"head"`
		Body     ttBody   `xml:"body"`
	}

	// Convert parser.ITTDocument to the marshalable ttRoot structure
	outputDoc := ttRoot{
		Xmlns:    "http://www.w3.org/ns/ttml",
		XmlnsTTP: "http://www.w3.org/ns/ttml#parameter",
		XmlnsTTS: "http://www.w3.org/ns/ttml#styling",
		TimeBase: "media", // As per GUIDE.md
		Lang:     doc.Lang,
	}

	// Styles
	for _, style := range doc.Styles {
		outputDoc.Head.Styling.Styles = append(outputDoc.Head.Styling.Styles, ttStyle{
			ID:         style.ID,
			Color:      style.Color,
			FontFamily: style.FontFamily,
			FontSize:   style.FontSize,
			FontStyle:  style.FontStyle,
			FontWeight: style.FontWeight,
		})
	}

	// Regions
	for _, region := range doc.Regions {
		outputDoc.Head.Layout.Regions = append(outputDoc.Head.Layout.Regions, ttRegion{
			ID:           region.ID,
			Origin:       region.Origin,
			Extent:       region.Extent,
			DisplayAlign: region.DisplayAlign,
			TextAlign:    region.TextAlign,
		})
	}

	// Cues
	var cues []ttP
	for _, cue := range doc.Cues {
		begin, err := formatTTMLTimestamp(cue.Begin)
		if err != nil {
			return "", err
		}
		end, err := formatTTMLTimestamp(cue.End)
		if err != nil {
			return "", err
		}
		cues = append(cues, ttP{
			Begin:   begin,
			End:     end,
			Content: cue.Content,
			Region:  cue.RegionID,
			Style:   strings.Join(cue.StyleIDs, " "),
		})
	}
	outputDoc.Body.Divs = append(outputDoc.Body.Divs, ttDiv{Ps: cues})

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	if err := encoder.Encode(outputDoc); err != nil {
		return "", err
	}

	// Perform a lightweight validation to ensure we didn't generate malformed XML.
	if err := ValidateTTML(buf.String()); err != nil {
		return "", fmt.Errorf("generated TTML failed validation: %w", err)
	}

	return buf.String(), nil
}

// formatTTMLTimestamp converts a big.Rat (in milliseconds) to a TTML time string (HH:MM:SS.ms).
func formatTTMLTimestamp(ms *big.Rat) (string, error) {
	if ms == nil {
		return "00:00:00.000", nil
	}
	msInt := new(big.Int).Quo(ms.Num(), ms.Denom()).Int64()

	hours := msInt / 3600000
	msInt %= 3600000
	minutes := msInt / 60000
	msInt %= 60000
	seconds := msInt / 1000
	milliseconds := msInt % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds), nil
}
