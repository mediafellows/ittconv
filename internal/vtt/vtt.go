package vtt

import (
	"encoding/xml"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/mediafellows/ittconv/internal/parser"
)

type textRun struct {
	text   string
	class  string
	italic bool
}

// ToVTT converts an ITTDocument to a WebVTT formatted string.
func ToVTT(doc *parser.ITTDocument) (string, error) {
	blocks := make([]string, 0, len(doc.Cues))
	for _, cue := range doc.Cues {
		lines, err := renderPayload(cue.Content, doc.Styles)
		if err != nil {
			return "", fmt.Errorf("rendering cue %q: %w", cue.ID, err)
		}

		blocks = append(blocks, fmt.Sprintf(
			"%s --> %s %s\r\n%s",
			formatTimestamp(cue.Begin),
			formatTimestamp(cue.End),
			cueSettings(cue, doc.Regions, len(lines)),
			strings.Join(lines, "\r\n"),
		))
	}

	return "\ufeffWEBVTT\r\n\r\n" + strings.Join(blocks, "\r\n\r\n"), nil
}

func cueSettings(cue parser.Cue, regions map[string]parser.Region, lineCount int) string {
	region := regions[cue.RegionID]
	align := region.TextAlign
	if align == "" {
		align = "center"
	}

	line := "-1"
	if region.DisplayAlign == "before" || cue.RegionID == "top" {
		line = "0"
	} else if lineCount > 1 {
		line = fmt.Sprintf("-%d", lineCount)
	}

	return fmt.Sprintf("align:%s line:%s", align, line)
}

func formatTimestamp(ms *big.Rat) string {
	msInt := roundMilliseconds(ms)

	hours := msInt / 3600000
	msInt %= 3600000
	minutes := msInt / 60000
	msInt %= 60000
	seconds := msInt / 1000
	milliseconds := msInt % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func roundMilliseconds(ms *big.Rat) int64 {
	if ms == nil || ms.Sign() <= 0 {
		return 0
	}

	rounded := new(big.Rat).Add(ms, big.NewRat(1, 2))
	return new(big.Int).Quo(rounded.Num(), rounded.Denom()).Int64()
}

func renderPayload(content string, styles map[string]parser.Style) ([]string, error) {
	lines, err := parseContent(content, styles)
	if err != nil {
		return nil, err
	}

	rendered := make([]string, 0, len(lines))
	for _, line := range lines {
		rendered = append(rendered, renderLine(line))
	}
	return rendered, nil
}

func parseContent(content string, styles map[string]parser.Style) ([][]textRun, error) {
	decoder := xml.NewDecoder(strings.NewReader("<root>" + content + "</root>"))
	styleStack := [][]string{{}}
	lines := [][]textRun{{}}

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "root":
				continue
			case "br":
				lines = append(lines, []textRun{})
			case "span":
				styleIDs := append([]string{}, styleStack[len(styleStack)-1]...)
				for _, attr := range t.Attr {
					if attr.Name.Local == "style" {
						styleIDs = append(styleIDs, strings.Fields(attr.Value)...)
					}
				}
				styleStack = append(styleStack, styleIDs)
			default:
				return nil, fmt.Errorf("unsupported content element <%s>", t.Name.Local)
			}
		case xml.EndElement:
			if t.Name.Local == "span" && len(styleStack) > 1 {
				styleStack = styleStack[:len(styleStack)-1]
			}
		case xml.CharData:
			if len(t) == 0 {
				continue
			}
			class, italic := resolveStyle(styleStack[len(styleStack)-1], styles)
			last := len(lines) - 1
			lines[last] = append(lines[last], textRun{
				text:   string(t),
				class:  class,
				italic: italic,
			})
		}
	}

	return lines, nil
}

func resolveStyle(styleIDs []string, styles map[string]parser.Style) (string, bool) {
	var class string
	var italic bool

	for _, id := range styleIDs {
		style, ok := styles[id]
		if !ok {
			continue
		}
		if style.FontStyle == "italic" {
			italic = true
		}
		if c := colorClass(style.Color); c != "" {
			class = c
		}
	}

	return class, italic
}

func colorClass(color string) string {
	switch strings.ToLower(strings.TrimSpace(color)) {
	case "", "white":
		return ""
	case "aqua":
		return "cyan"
	case "yellow":
		return "yellow"
	case "cyan":
		return "cyan"
	}

	if strings.HasPrefix(color, "#") {
		return "color" + strings.ToUpper(strings.TrimPrefix(color, "#"))
	}

	return color
}

func renderLine(runs []textRun) string {
	runs = mergeRuns(runs)

	var b strings.Builder
	for i := 0; i < len(runs); {
		run := runs[i]
		if run.italic {
			end, hasDifferentClass := italicGroup(runs, i)
			if hasDifferentClass {
				b.WriteString("<i>")
				for i < end {
					writeRun(&b, runs[i].text, runs[i].class, false)
					i++
				}
				b.WriteString("</i>")
				continue
			}
		}

		if run.class == "" {
			writeRun(&b, run.text, "", run.italic)
			i++
			continue
		}

		b.WriteString("<c.")
		b.WriteString(run.class)
		b.WriteByte('>')
		for i < len(runs) && runs[i].class == run.class {
			writeRun(&b, runs[i].text, "", runs[i].italic)
			i++
		}
		b.WriteString("</c>")
	}

	return b.String()
}

func italicGroup(runs []textRun, start int) (int, bool) {
	class := runs[start].class
	for i := start + 1; i < len(runs); i++ {
		if !runs[i].italic {
			return i, false
		}
		if runs[i].class != class {
			for j := i + 1; j < len(runs) && runs[j].italic; j++ {
				i = j
			}
			return i + 1, true
		}
	}
	return len(runs), false
}

func mergeRuns(runs []textRun) []textRun {
	merged := make([]textRun, 0, len(runs))
	for _, run := range runs {
		if run.text == "" {
			continue
		}
		last := len(merged) - 1
		if last >= 0 && merged[last].class == run.class && merged[last].italic == run.italic {
			merged[last].text += run.text
			continue
		}
		merged = append(merged, run)
	}
	return merged
}

func writeRun(b *strings.Builder, text, class string, italic bool) {
	if class != "" {
		b.WriteString("<c.")
		b.WriteString(class)
		b.WriteByte('>')
	}
	if italic {
		b.WriteString("<i>")
	}

	b.WriteString(escapeText(text))

	if italic {
		b.WriteString("</i>")
	}
	if class != "" {
		b.WriteString("</c>")
	}
}

func escapeText(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	return strings.ReplaceAll(text, "<", "&lt;")
}
