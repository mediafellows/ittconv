package itt2vtt

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/orisano/gosax"
)

const (
	webvttHeader = "WEBVTT"
)

type Style struct {
	ID              string
	FontWeight      string
	FontStyle       string
	TextDecoration  string
	Color           string
	BackgroundColor string
}

type Region struct {
	ID     string
	Origin string
	Extent string
	Align  string
}

type Cue struct {
	ID       string
	Start    string
	End      string
	Settings string
	Text     string
}

type Converter struct {
	frameRate  float64
	styles     map[string]Style
	regions    map[string]Region
	cues       []Cue
	lang       string
	currentCue *Cue
	styleStack []Style
	cssStyles  map[string]string
}

func NewConverter() *Converter {
	return &Converter{
		frameRate: 30, // Default frame rate
		styles:    make(map[string]Style),
		regions:   make(map[string]Region),
		cssStyles: make(map[string]string),
	}
}

func (c *Converter) Convert(r io.Reader) (string, error) {
	reader := gosax.NewReader(r)
	for {
		e, err := reader.Event()
		if e.Type() == gosax.EventEOF {
			break
		}
		if err != nil {
			return "", err
		}

		switch e.Type() {
		case gosax.EventStart:
			se, err := gosax.StartElement(e.Bytes)
			if err != nil {
				return "", err
			}
			c.handleStartElement(se)
		case gosax.EventText:
			cd, err := gosax.CharData(e.Bytes)
			if err != nil {
				return "", err
			}
			c.handleCharData(cd)
		case gosax.EventEnd:
			ee := gosax.EndElement(e.Bytes)
			c.handleEndElement(ee)
		}
	}

	var out bytes.Buffer
	out.WriteString(webvttHeader)
	out.WriteString("\r\n\r\n")

	if len(c.cssStyles) > 0 {
		out.WriteString("STYLE\r\n")
		for class, rule := range c.cssStyles {
			out.WriteString(fmt.Sprintf("::cue(.%s) { %s }\r\n", class, rule))
		}
		out.WriteString("\r\n")
	}

	if c.lang != "" {
		out.WriteString("Language: ")
		out.WriteString(c.lang)
		out.WriteString("\r\n\r\n")
	}

	for _, cue := range c.cues {
		out.WriteString(fmt.Sprintf("%s --> %s %s\r\n", cue.Start, cue.End, cue.Settings))
		out.WriteString(cue.Text)
		out.WriteString("\r\n\r\n")
	}

	return out.String(), nil
}

func (c *Converter) handleStartElement(se xml.StartElement) {
	switch se.Name.Local {
	case "tt":
		for _, attr := range se.Attr {
			if attr.Name.Local == "lang" {
				c.lang = attr.Value
			}
			if attr.Name.Space == "ttp" && attr.Name.Local == "frameRate" {
				if fr, err := strconv.ParseFloat(attr.Value, 64); err == nil {
					c.frameRate = fr
				}
			}
		}
	case "style":
		var s Style
		s.ID = getAttr(se.Attr, "id")
		s.FontWeight = getAttr(se.Attr, "fontWeight")
		s.FontStyle = getAttr(se.Attr, "fontStyle")
		s.TextDecoration = getAttr(se.Attr, "textDecoration")
		s.Color = getAttr(se.Attr, "color")
		s.BackgroundColor = getAttr(se.Attr, "backgroundColor")
		if s.ID != "" {
			c.styles[s.ID] = s
			var rules []string
			if s.Color != "" {
				rules = append(rules, "color: "+s.Color)
			}
			if s.BackgroundColor != "" {
				rules = append(rules, "background-color: "+s.BackgroundColor)
			}
			if len(rules) > 0 {
				c.cssStyles[s.ID] = strings.Join(rules, "; ")
			}
		}
	case "region":
		var r Region
		r.ID = getAttr(se.Attr, "id")
		r.Origin = getAttr(se.Attr, "origin")
		r.Extent = getAttr(se.Attr, "extent")
		r.Align = getAttr(se.Attr, "textAlign")
		if r.ID != "" {
			c.regions[r.ID] = r
		}
	case "p":
		c.currentCue = &Cue{}
		c.currentCue.Start = c.toTime(getAttr(se.Attr, "begin"))
		c.currentCue.End = c.toTime(getAttr(se.Attr, "end"))
		if regionID := getAttr(se.Attr, "region"); regionID != "" {
			if region, ok := c.regions[regionID]; ok {
				c.currentCue.Settings = c.regionToSettings(region)
			}
		}
		// NOTE: p can have style attribute as well.
		// For simplicity, we only handle span styles for now.
	case "span":
		if styleID := getAttr(se.Attr, "style"); styleID != "" {
			if style, ok := c.styles[styleID]; ok {
				c.currentCue.Text += c.styleToTags(style, true)
				c.styleStack = append(c.styleStack, style)
			}
		}
	case "br":
		if c.currentCue != nil {
			c.currentCue.Text += "\n"
		}
	}
}

func (c *Converter) handleEndElement(ee xml.EndElement) {
	switch ee.Name.Local {
	case "p":
		if c.currentCue != nil {
			if c.currentCue.Settings == "" {
				// Default region
				c.currentCue.Settings = "align:center position:50% line:90%"
			}

			// Clean up whitespace
			c.currentCue.Text = strings.TrimSpace(c.currentCue.Text)

			c.cues = append(c.cues, *c.currentCue)
			c.currentCue = nil
		}
	case "span":
		if len(c.styleStack) > 0 {
			style := c.styleStack[len(c.styleStack)-1]
			c.styleStack = c.styleStack[:len(c.styleStack)-1]
			c.currentCue.Text += c.styleToTags(style, false)
		}
	}
}

func (c *Converter) handleCharData(cd xml.CharData) {
	if c.currentCue != nil {
		c.currentCue.Text += string(cd)
	}
}

func (c *Converter) toTime(ittTime string) string {
	if ittTime == "" {
		return ""
	}
	// HH:MM:SS:FF or HH:MM:SS.mmm
	parts := strings.Split(ittTime, ":")
	if len(parts) == 4 { // Frame based
		h, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		s, _ := strconv.Atoi(parts[2])
		f, _ := strconv.Atoi(parts[3])
		ms := (float64(h*3600+m*60+s) * 1000) + (float64(f*1000) / c.frameRate)
		h = int(ms / 3600000)
		ms -= float64(h * 3600000)
		m = int(ms / 60000)
		ms -= float64(m * 60000)
		s = int(ms / 1000)
		ms -= float64(s * 1000)
		return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, int(ms))
	} else if strings.Contains(ittTime, ".") { // Time based with ms
		parts := strings.Split(ittTime, ".")
		if len(parts) == 2 && len(parts[1]) > 3 {
			parts[1] = parts[1][:3]
			return strings.Join(parts, ".")
		}
		return ittTime
	} else { // Time based without ms
		return ittTime + ".000"
	}
}

var reOrigin = regexp.MustCompile(`(\d+(\.\d+)?)%\s+(\d+(\.\d+)?)%`)
var reExtent = regexp.MustCompile(`(\d+(\.\d+)?)%\s+(\d+(\.\d+)?)%`)

func (c *Converter) regionToSettings(r Region) string {
	var settings []string
	if matches := reOrigin.FindStringSubmatch(r.Origin); len(matches) > 3 {
		settings = append(settings, "line:"+matches[3]+"%")
	}
	if matches := reExtent.FindStringSubmatch(r.Extent); len(matches) > 0 {
		settings = append(settings, "size:"+matches[1]+"%")
	}
	if r.Align != "" {
		settings = append(settings, "align:"+r.Align)
	}
	// Default position
	settings = append(settings, "position:50%")
	return strings.Join(settings, " ")
}

func (c *Converter) styleToTags(s Style, opening bool) string {
	var tags string
	if opening {
		if s.Color != "" || s.BackgroundColor != "" {
			tags += "<c." + s.ID + ">"
		}
		if s.FontWeight == "bold" {
			tags += "<b>"
		}
		if s.FontStyle == "italic" {
			tags += "<i>"
		}
		if s.TextDecoration == "underline" {
			tags += "<u>"
		}
	} else {
		if s.TextDecoration == "underline" {
			tags += "</u>"
		}
		if s.FontStyle == "italic" {
			tags += "</i>"
		}
		if s.FontWeight == "bold" {
			tags += "</b>"
		}
		if s.Color != "" || s.BackgroundColor != "" {
			tags += "</c>"
		}
	}
	return tags
}

func getAttr(attrs []xml.Attr, name string) string {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}
