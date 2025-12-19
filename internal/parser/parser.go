package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/mediafellows/ittconv/internal/timecode"

	"github.com/orisano/gosax"
)

var logger *slog.Logger

func init() {
	logLevel := slog.LevelInfo
	if os.Getenv("DEBUG") == "1" {
		logLevel = slog.LevelDebug
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
}

// ParseITT parses an ITT XML string into an ITTDocument structure.
func ParseITT(ittSource string) (*ITTDocument, error) {
	logger.Debug("Starting ITT parsing")
	doc := &ITTDocument{
		Styles:  make(map[string]Style),
		Regions: make(map[string]Region),
		Cues:    []Cue{},
	}

	reader := strings.NewReader(ittSource)
	r := gosax.NewReader(reader)
	r.EmitSelfClosingTag = true // Ensure self-closing tags are recognized

	handler := &ittHandler{doc: doc, reader: r} // Pass reader to handler

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
			logger.Debug("Handling start element", "name", startElement.Name.Local)
			if err := handler.handleStartElement(startElement.Name, startElement.Attr); err != nil {
				return nil, err
			}
		case gosax.EventEnd:
			endElement := gosax.EndElement(e.Bytes)
			logger.Debug("Handling end element", "name", endElement.Name.Local)
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
	fr := doc.FrameRateValue
	if fr == nil {
		baseFrameRate, err := timecode.NewFrameRate(doc.FrameRate)
		if err != nil {
			return nil, fmt.Errorf("invalid frame rate '%s': %w", doc.FrameRate, err)
		}
		if doc.FrameRateMultiplierNum > 0 && doc.FrameRateMultiplierDen > 0 {
			if baseFrameRate.IsInt() {
				baseFrameRate = &timecode.FrameRate{
					Rat: new(big.Rat).Mul(baseFrameRate.Rat, big.NewRat(int64(doc.FrameRateMultiplierNum), int64(doc.FrameRateMultiplierDen))),
				}
			} else {
				logger.Debug("Skipping frameRateMultiplier because base frameRate is non-integer", "frameRate", doc.FrameRate)
			}
		}
		fr = baseFrameRate
		doc.FrameRateValue = fr
	}
	logger.Debug("Successfully parsed framerate", "framerate", doc.FrameRate)

	for i := range doc.Cues {
		cue := &doc.Cues[i]
		if cue.BeginTimecode != nil {
			ms, err := cue.BeginTimecode.ToMilliseconds(fr)
			if err != nil {
				return nil, fmt.Errorf("error converting begin timecode '%v': %w", cue.BeginTimecode, err)
			}
			cue.Begin = ms
			logger.Debug("Converted begin timecode", "smpte", cue.BeginTimecode, "ms", ms)
		}
		if cue.EndTimecode != nil {
			ms, err := cue.EndTimecode.ToMilliseconds(fr)
			if err != nil {
				return nil, fmt.Errorf("error converting end timecode '%v': %w", cue.EndTimecode, err)
			}
			cue.End = ms
			logger.Debug("Converted end timecode", "smpte", cue.EndTimecode, "ms", ms)
		}

		if cue.Offset != nil {
			if cue.Begin != nil {
				cue.Begin = new(big.Rat).Add(cue.Begin, cue.Offset)
			}
			if cue.End != nil {
				cue.End = new(big.Rat).Add(cue.End, cue.Offset)
			}
		}

		if cue.Begin != nil && cue.Begin.Sign() < 0 {
			logger.Warn("Clamping negative cue begin", "id", cue.ID, "value", cue.Begin.String())
			cue.Begin = big.NewRat(0, 1)
		}
		if cue.End != nil && cue.End.Sign() < 0 {
			logger.Warn("Clamping negative cue end", "id", cue.ID, "value", cue.End.String())
			cue.End = big.NewRat(0, 1)
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
	reader        *gosax.Reader
	offsetStack   []*big.Rat
	frameRate     *timecode.FrameRate
}

func (h *ittHandler) handleStartElement(name xml.Name, attrs []xml.Attr) error {
	if h.inPElement {
		// If we are inside a <p> element, treat everything as raw content
		var buf bytes.Buffer
		buf.WriteByte('<')
		buf.WriteString(name.Local)
		for _, attr := range attrs {
			buf.WriteByte(' ')
			buf.WriteString(attr.Name.Local)
			buf.WriteString(`="`)
			buf.WriteString(attr.Value)
			buf.WriteByte('"')
		}

		// Handle self-closing tags like <br/>
		if name.Local == "br" {
			buf.WriteString("/>")
		} else {
			buf.WriteByte('>')
		}

		h.contentBuffer.Write(buf.Bytes())

		return nil
	}

	switch name.Local {
	case "tt":
		var frameRateMultiplier string
		for _, attr := range attrs {
			switch attr.Name.Local {
			case "lang":
				h.doc.Lang = attr.Value
				logger.Debug("Parsed lang", "value", attr.Value)
			case "timeBase":
				h.doc.TimeBase = attr.Value
				logger.Debug("Parsed timeBase", "value", attr.Value)
			case "frameRate":
				h.doc.FrameRate = attr.Value
				logger.Debug("Parsed frameRate", "value", attr.Value)
			case "frameRateMultiplier":
				frameRateMultiplier = attr.Value
				logger.Debug("Parsed frameRateMultiplier", "value", attr.Value)
			}
		}

		if frameRateMultiplier != "" {
			parts := strings.Fields(frameRateMultiplier)
			if len(parts) != 2 {
				return fmt.Errorf("invalid frameRateMultiplier format: %s", frameRateMultiplier)
			}
			num, err := strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("invalid frameRateMultiplier numerator %q: %w", parts[0], err)
			}
			den, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid frameRateMultiplier denominator %q: %w", parts[1], err)
			}
			if num <= 0 || den <= 0 {
				return fmt.Errorf("frameRateMultiplier values must be positive: %s", frameRateMultiplier)
			}
			h.doc.FrameRateMultiplierNum = num
			h.doc.FrameRateMultiplierDen = den
		}

		if h.doc.FrameRate != "" {
			fr, err := timecode.NewFrameRate(h.doc.FrameRate)
			if err != nil {
				return fmt.Errorf("invalid frame rate '%s': %w", h.doc.FrameRate, err)
			}
			if h.doc.FrameRateMultiplierNum > 0 && h.doc.FrameRateMultiplierDen > 0 {
				if fr.IsInt() {
					fr = &timecode.FrameRate{
						Rat: new(big.Rat).Mul(fr.Rat, big.NewRat(int64(h.doc.FrameRateMultiplierNum), int64(h.doc.FrameRateMultiplierDen))),
					}
				} else {
					logger.Debug("Skipping frameRateMultiplier because base frameRate is non-integer", "frameRate", h.doc.FrameRate)
				}
			}
			h.frameRate = fr
			h.doc.FrameRateValue = fr
			logger.Debug("Computed effective framerate", "value", fr.String())
		}
	case "body", "div":
		regionFromAttr := ""
		var beginAttr string
		for _, attr := range attrs {
			if attr.Name.Local == "region" {
				regionFromAttr = attr.Value
			}
			if attr.Name.Local == "begin" {
				beginAttr = attr.Value
			}
		}
		h.regionStack = append(h.regionStack, regionFromAttr)
		if err := h.pushOffset(beginAttr, name.Local); err != nil {
			return err
		}
	case "style":
		h.currentStyle = &Style{}
		for _, attr := range attrs {
			switch attr.Name.Local {
			case "id":
				h.currentStyle.ID = attr.Value
				logger.Debug("Parsed style id", "value", attr.Value)
			case "fontFamily":
				h.currentStyle.FontFamily = attr.Value
				logger.Debug("Parsed style fontFamily", "value", attr.Value)
			case "fontSize":
				h.currentStyle.FontSize = attr.Value
				logger.Debug("Parsed style fontSize", "value", attr.Value)
			case "fontWeight":
				h.currentStyle.FontWeight = attr.Value
				logger.Debug("Parsed style fontWeight", "value", attr.Value)
			case "fontStyle":
				h.currentStyle.FontStyle = attr.Value
				logger.Debug("Parsed style fontStyle", "value", attr.Value)
			case "color":
				h.currentStyle.Color = attr.Value
				logger.Debug("Parsed style color", "value", attr.Value)
			}
		}
		if h.currentStyle.ID != "" {
			h.doc.Styles[h.currentStyle.ID] = *h.currentStyle
			logger.Debug("Stored style", "id", h.currentStyle.ID, "details", *h.currentStyle)
		}
	case "region":
		h.currentRegion = &Region{}
		for _, attr := range attrs {
			switch attr.Name.Local {
			case "id":
				h.currentRegion.ID = attr.Value
				logger.Debug("Parsed region id", "value", attr.Value)
			case "origin":
				h.currentRegion.Origin = attr.Value
				logger.Debug("Parsed region origin", "value", attr.Value)
			case "extent":
				h.currentRegion.Extent = attr.Value
				logger.Debug("Parsed region extent", "value", attr.Value)
			case "textAlign":
				h.currentRegion.TextAlign = attr.Value
				logger.Debug("Parsed region textAlign", "value", attr.Value)
			case "displayAlign":
				h.currentRegion.DisplayAlign = attr.Value
				logger.Debug("Parsed region displayAlign", "value", attr.Value)
			}
		}
		if h.currentRegion.ID != "" {
			h.doc.Regions[h.currentRegion.ID] = *h.currentRegion
			logger.Debug("Stored region", "id", h.currentRegion.ID, "details", *h.currentRegion)
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
					logger.Warn("Invalid begin timecode format", "value", attr.Value, "error", err)
					// Continue parsing, but cue.BeginTimecode will be nil, handled later
				} else {
					h.currentCue.BeginTimecode = tc
				}
				h.currentCue.ID = attr.Value // For now, use begin time as ID
			case "end":
				tc, err := timecode.ParseSMPTETimecode(attr.Value)
				if err != nil {
					logger.Warn("Invalid end timecode format", "value", attr.Value, "error", err)
				} else {
					h.currentCue.EndTimecode = tc
				}
			case "region":
				pRegion = attr.Value
				hasPRegion = true
				logger.Debug("Parsed p region", "value", attr.Value)
			case "style":
				h.currentCue.StyleIDs = strings.Fields(attr.Value)
				logger.Debug("Parsed p style", "value", attr.Value)
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
		h.currentCue.Offset = h.currentOffset()
		logger.Debug("Starting p element", "id", h.currentCue.ID, "region", h.currentCue.RegionID)
	case "span":
		h.inSpanElement = true
		// Handle span styles if needed, add to currentCue.StyleIDs
		for _, attr := range attrs {
			if attr.Name.Local == "style" {
				h.currentCue.StyleIDs = append(h.currentCue.StyleIDs, strings.Fields(attr.Value)...)
				logger.Debug("Parsed span style", "value", attr.Value)
			}
		}
	case "br":
		if h.inPElement || h.inSpanElement {
			h.contentBuffer.WriteString("\n")
		}
	}
	return nil
}

func (h *ittHandler) pushOffset(beginAttr string, elementName string) error {
	var parent *big.Rat
	if len(h.offsetStack) > 0 && h.offsetStack[len(h.offsetStack)-1] != nil {
		parent = new(big.Rat).Set(h.offsetStack[len(h.offsetStack)-1])
	}

	offset := parent
	if beginAttr != "" {
		if h.frameRate == nil {
			return fmt.Errorf("frameRate attribute missing in <tt> tag")
		}
		tc, err := timecode.ParseSMPTETimecode(beginAttr)
		if err != nil {
			return fmt.Errorf("error parsing begin timecode '%s' on <%s>: %w", beginAttr, elementName, err)
		}
		ms, err := tc.ToMilliseconds(h.frameRate)
		if err != nil {
			return fmt.Errorf("error converting begin timecode '%s' on <%s>: %w", beginAttr, elementName, err)
		}
		if offset == nil {
			offset = new(big.Rat)
		}
		offset.Add(offset, ms)
	}

	h.offsetStack = append(h.offsetStack, offset)
	return nil
}

func (h *ittHandler) currentOffset() *big.Rat {
	if len(h.offsetStack) == 0 || h.offsetStack[len(h.offsetStack)-1] == nil {
		return nil
	}
	return new(big.Rat).Set(h.offsetStack[len(h.offsetStack)-1])
}

func (h *ittHandler) handleEndElement(name xml.Name) error {
	if name.Local == "p" {
		if h.currentCue != nil {
			h.currentCue.Content = h.contentBuffer.String()
			h.doc.Cues = append(h.doc.Cues, *h.currentCue)
			logger.Debug("Finalized cue", "id", h.currentCue.ID, "content", h.currentCue.Content)
		}
		h.inPElement = false
		h.currentCue = nil
		h.contentBuffer.Reset()
		return nil
	}

	if h.inPElement {
		// Don't write a closing tag for self-closing tags
		if name.Local != "br" {
			var buf bytes.Buffer
			buf.WriteString("</")
			buf.WriteString(name.Local)
			buf.WriteString(">")
			h.contentBuffer.Write(buf.Bytes())
		}
		return nil
	}

	switch name.Local {
	case "span":
		h.inSpanElement = false
	case "body", "div":
		if len(h.regionStack) > 0 {
			h.regionStack = h.regionStack[:len(h.regionStack)-1]
		}
		if len(h.offsetStack) > 0 {
			h.offsetStack = h.offsetStack[:len(h.offsetStack)-1]
		}
	}
	return nil
}

func (h *ittHandler) handleCharData(c xml.CharData) error {
	if h.inPElement || h.inSpanElement {
		if err := xml.EscapeText(&h.contentBuffer, c); err != nil {
			return err
		}
	}
	return nil
}
