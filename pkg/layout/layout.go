package layout

import (
	"math"
	"strings"

	"fyne.io/fyne/v2"

	"github.com/kuosandys/browser-engineering/pkg/parser"
)

const (
	DefaultHStep   float32 = 13
	DefaultVStep   float32 = 16
	defaultLeading float32 = 1.25
)

type font struct {
	Style fyne.TextStyle
	Size  float32
}

type Layout struct {
	DisplayList []DisplayItem
	cursorX     float32
	cursorY     float32
	font        font
	leading     float32
	width       float32
	scale       float32
	line        []DisplayItem
}

type DisplayItem struct {
	X    float32
	Y    float32
	Text string
	Font font
}

func NewLayout(tokens []interface{}, width float32, scale float32) *Layout {
	l := &Layout{
		DisplayList: make([]DisplayItem, 0),
		cursorX:     DefaultHStep,
		cursorY:     0,
		font:        font{Style: fyne.TextStyle{}, Size: DefaultVStep * scale},
		leading:     defaultLeading,
		width:       width,
		scale:       scale,
	}

	l.createLayout(tokens)

	return l
}

func (l *Layout) createLayout(tokens []interface{}) {
	for _, tok := range tokens {
		l.token(tok)
	}

	l.flush()
}

func (l *Layout) token(token interface{}) {
	switch token.(type) {
	case *parser.Text:
		l.text(token.(*parser.Text))
	case *parser.Tag:
		switch token.(*parser.Tag).Tag {
		case "i":
			l.font.Style.Italic = true
		case "/i":
			l.font.Style.Italic = false
		case "b":
			l.font.Style.Bold = true
		case "/b":
			l.font.Style.Bold = false
		case "small":
			l.font.Size -= 2
		case "/small":
			l.font.Size += 2
		case "big":
			l.font.Size += 4
		case "/big":
			l.font.Size -= 4
		case "br":
			l.flush()
		case "/p":
			l.flush()
			l.cursorY += DefaultVStep * l.scale
		case "sub":
			l.font.Size /= 2
			l.leading *= 1.25
		case "/sub":
			l.font.Size *= 2
			l.leading /= 1.25
		case "sup":
			l.font.Size /= 2
			l.leading *= 2
		case "/sup":
			l.font.Size *= 2
			l.leading /= 2
		}
	}
}

func (l *Layout) text(token *parser.Text) {
	spaceSize := fyne.MeasureText(" ", l.font.Size, l.font.Style)

	for _, word := range strings.Fields(token.Text) {
		size := fyne.MeasureText(word, l.font.Size, l.font.Style)
		if l.cursorX+size.Width+spaceSize.Width >= l.width-(DefaultHStep) {
			l.flush()
		}

		l.cursorX += spaceSize.Width
		l.line = append(l.line, DisplayItem{X: l.cursorX, Y: size.Height * l.leading, Text: word, Font: l.font})
		l.cursorX += size.Width
	}
}

func (l *Layout) flush() {
	if len(l.line) <= 0 {
		return
	}

	var maxHeight float32
	for _, item := range l.line {
		maxHeight = float32(math.Max(float64(maxHeight), float64(fyne.MeasureText(item.Text, item.Font.Size, item.Font.Style).Height)))
	}
	baseline := l.cursorY + (maxHeight)
	for _, item := range l.line {
		l.DisplayList = append(l.DisplayList, DisplayItem{X: item.X, Y: baseline - item.Y, Text: item.Text, Font: item.Font})
	}

	l.line = make([]DisplayItem, 0)
	l.cursorY += maxHeight * l.leading
	l.cursorX = DefaultHStep
}
