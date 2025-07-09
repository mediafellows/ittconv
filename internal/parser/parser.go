package parser

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	"ittconv/internal/timecode"

	"github.com/orisano/gosax"
)

// ParseITT parses an ITT XML string into an ITTDocument structure.
func ParseITT(ittSource string) (*ITTDocument, error) {
	doc := &ITTDocument{
		Styles:  make(map[string]Style),
		Regions: make(map[string]Region),
		Cues:    []Cue{},
	}

	reader := strings.NewReader(ittSource)
	r := gosax.NewReader(reader)
	r.EmitSelfClosingTag = true // Ensure self-closing tags are recognized

	handler := &ittHandler{doc: doc} // Still use ittHandler to hold state and logic

	for {
		e, err := r.Event()
		if e.Type() == gosax.EventEOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading XML event: %w", err)
		}

		switch e.Type() {
		case gosax.EventStart:
			startElement, err := gosax.StartElement(e.Bytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing start element: %w", err)
			}
			if err := handler.handleStartElement(startElement.Name, startElement.Attr); err != nil {
				return nil, err
			}
		case gosax.EventEnd:
			endElement := gosax.EndElement(e.Bytes)
			if err := handler.handleEndElement(endElement.Name); err != nil {
				return nil, err
			}
		case gosax.EventText:
			charData, err := gosax.CharData(e.Bytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing character data: %w", err)
			}
			if err := handler.handleCharData(charData); err != nil {
				return nil, err
			}
			// Add other event types if needed (e.g., comments, processing instructions)
		}
	}

	// Post-processing: Convert SMPTE timecodes to milliseconds
	if doc.FrameRate == "" {
		return nil, fmt.Errorf("frameRate attribute missing in <tt> tag")
	}
	fr, err := timecode.NewFrameRate(doc.FrameRate)
	if err != nil {
		return nil, fmt.Errorf("invalid frame rate '%s': %w", doc.FrameRate, err)
	}

	for i := range doc.Cues {
		cue := &doc.Cues[i]
		if cue.BeginTimecode != nil {
			ms, err := cue.BeginTimecode.ToMilliseconds(fr)
			if err != nil {
				return nil, fmt.Errorf("error converting begin timecode '%v': %w", cue.BeginTimecode, err)
			}
			cue.Begin = ms
		}
		if cue.EndTimecode != nil {
			ms, err := cue.EndTimecode.ToMilliseconds(fr)
			if err != nil {
				return nil, fmt.Errorf("error converting end timecode '%v': %w", cue.EndTimecode, err)
			}
			cue.End = ms
		}

		// Validate begin < end
		if cue.Begin != nil && cue.End != nil && cue.Begin.Cmp(cue.End) >= 0 {
			return nil, fmt.Errorf("invalid cue timing: begin time (%s) is not less than end time (%s) for cue ID %s",
				cue.Begin.String(), cue.End.String(), cue.ID)
		}
	}

	return doc, nil
}

type ittHandler struct {
	doc           *ITTDocument
	currentCue    *Cue
	currentStyle  *Style
	currentRegion *Region
	contentBuffer strings.Builder
	inPElement    bool
	inSpanElement bool
	regionStack   []string
}

func (h *ittHandler) handleStartElement(name xml.Name, attrs []xml.Attr) error {
	switch name.Local {
	case "tt":
		for _, attr := range attrs {
			switch attr.Name.Local {
			case "lang":
				h.doc.Lang = attr.Value
			case "timeBase":
				h.doc.TimeBase = attr.Value
			case "frameRate":
				h.doc.FrameRate = attr.Value
			}
		}
	case "body", "div":
		regionFromAttr := ""
		for _, attr := range attrs {
			if attr.Name.Local == "region" {
				regionFromAttr = attr.Value
			}
		}
		h.regionStack = append(h.regionStack, regionFromAttr)
	case "style":
		h.currentStyle = &Style{}
		for _, attr := range attrs {
			switch attr.Name.Local {
			case "id":
				h.currentStyle.ID = attr.Value
			case "fontFamily":
				h.currentStyle.FontFamily = attr.Value
			case "fontSize":
				h.currentStyle.FontSize = attr.Value
			case "fontWeight":
				h.currentStyle.FontWeight = attr.Value
			case "fontStyle":
				h.currentStyle.FontStyle = attr.Value
			case "color":
				h.currentStyle.Color = attr.Value
			}
		}
		if h.currentStyle.ID != "" {
			h.doc.Styles[h.currentStyle.ID] = *h.currentStyle
		}
	case "region":
		h.currentRegion = &Region{}
		for _, attr := range attrs {
			switch attr.Name.Local {
			case "id":
				h.currentRegion.ID = attr.Value
			case "origin":
				h.currentRegion.Origin = attr.Value
			case "extent":
				h.currentRegion.Extent = attr.Value
			case "textAlign":
				h.currentRegion.TextAlign = attr.Value
			case "displayAlign":
				h.currentRegion.DisplayAlign = attr.Value
			}
		}
		if h.currentRegion.ID != "" {
			h.doc.Regions[h.currentRegion.ID] = *h.currentRegion
		}
	case "p":
		h.inPElement = true
		h.currentCue = &Cue{}
		var pRegion string
		var hasPRegion bool

		for _, attr := range attrs {
			switch attr.Name.Local {
			case "begin":
				tc, err := timecode.ParseSMPTETimecode(attr.Value)
				if err != nil {
					log.Printf("Warning: invalid begin timecode format '%s': %v", attr.Value, err)
					// Continue parsing, but cue.BeginTimecode will be nil, handled later
				} else {
					h.currentCue.BeginTimecode = tc
				}
				h.currentCue.ID = attr.Value // For now, use begin time as ID
			case "end":
				tc, err := timecode.ParseSMPTETimecode(attr.Value)
				if err != nil {
					log.Printf("Warning: invalid end timecode format '%s': %v", attr.Value, err)
				} else {
					h.currentCue.EndTimecode = tc
				}
			case "region":
				pRegion = attr.Value
				hasPRegion = true
			case "style":
				h.currentCue.StyleIDs = strings.Fields(attr.Value)
			}
		}

		if hasPRegion {
			h.currentCue.RegionID = pRegion
		} else {
			// Inherit from stack, finding the closest ancestor with a region
			for i := len(h.regionStack) - 1; i >= 0; i-- {
				if h.regionStack[i] != "" {
					h.currentCue.RegionID = h.regionStack[i]
					break
				}
			}
		}
	case "span":
		h.inSpanElement = true
		// Handle span styles if needed, add to currentCue.StyleIDs
		for _, attr := range attrs {
			if attr.Name.Local == "style" {
				h.currentCue.StyleIDs = append(h.currentCue.StyleIDs, strings.Fields(attr.Value)...)
			}
		}
	case "br":
		if h.inPElement || h.inSpanElement {
			h.contentBuffer.WriteString("\n")
		}
	}
	return nil
}

func (h *ittHandler) handleEndElement(name xml.Name) error {
	switch name.Local {
	case "p":
		if h.currentCue != nil {
			h.currentCue.Content = h.contentBuffer.String()
			h.doc.Cues = append(h.doc.Cues, *h.currentCue)
		}
		h.inPElement = false
		h.currentCue = nil
		h.contentBuffer.Reset()
	case "span":
		h.inSpanElement = false
	case "body", "div":
		if len(h.regionStack) > 0 {
			h.regionStack = h.regionStack[:len(h.regionStack)-1]
		}
	}
	return nil
}

func (h *ittHandler) handleCharData(c xml.CharData) error {
	// log.Printf("CharData: %s", c) // This can be very verbose
	if h.inPElement || h.inSpanElement {
		h.contentBuffer.Write([]byte(c))
	}
	return nil
}
