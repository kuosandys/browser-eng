package layout

import (
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"

	"github.com/kuosandys/browser-engineering/pkg/parser"
)

const (
	DefaultHStep   float32 = 13
	DefaultVStep   float32 = 16
	defaultLeading float32 = 1.25
)

type Layout struct {
	DisplayList []DisplayItem
	cursorX     float32
	cursorY     float32
	fontStyle   fyne.TextStyle
	fontSize    float32
	width       float32
	scale       float32
	line        []DisplayItem
}

type DisplayItem struct {
	X         float32
	Y         float32
	Text      string
	FontStyle fyne.TextStyle
	FontSize  float32
}

func NewLayout(tokens []interface{}, width float32, scale float32) *Layout {
	l := &Layout{
		DisplayList: make([]DisplayItem, 0),
		cursorX:     DefaultHStep,
		cursorY:     0,
		fontStyle:   fyne.TextStyle{},
		fontSize:    DefaultVStep * scale,
		width:       width,
		scale:       scale,
	}

	l.createLayout(tokens)

	return l
}

func (l *Layout) token(token interface{}, inEmoji *bool) {
	switch token.(type) {
	case *parser.Text:
		l.text(token.(*parser.Text), inEmoji)
	case *parser.Tag:
		switch token.(*parser.Tag).Tag {
		case "i":
			l.fontStyle.Italic = true
		case "/i":
			l.fontStyle.Italic = false
		case "b":
			l.fontStyle.Bold = true
		case "/b":
			l.fontStyle.Bold = false
		case "small":
			l.fontSize -= 2
		case "/small":
			l.fontSize += 2
		case "big":
			l.fontSize += 4
		case "/big":
			l.fontSize -= 4
		case "br":
			l.flush()
		case "/p":
			l.flush()
			l.cursorY += DefaultVStep * l.scale
		}
	}
}

func (l *Layout) text(token *parser.Text, inEmoji *bool) {
	spaceSize := fyne.MeasureText(" ", l.fontSize, l.fontStyle)

	for _, word := range strings.Fields(token.Text) {
		size := fyne.MeasureText(word, l.fontSize, l.fontStyle)
		l.line = append(l.line, DisplayItem{X: l.cursorX, Y: size.Height, Text: word, FontStyle: l.fontStyle, FontSize: l.fontSize})

		ascii := strings.Trim(strconv.QuoteToASCII(word), "\"")
		if len(ascii) > 2 && (ascii[0:2] == "\\U" || ascii[0:2] == "\\u") {
			// don't change cursor position for emoji unicode characters
			*inEmoji = true
			continue
		}

		hStep := DefaultHStep * l.scale
		if *inEmoji {
			// use two HSteps for emoji unicode characters
			hStep *= 2
		}

		if l.cursorX+size.Width+spaceSize.Width >= l.width-(5*DefaultHStep) {
			l.flush()
		} else {
			l.cursorX += size.Width + spaceSize.Width
		}
	}
}

func (l *Layout) flush() {
	if len(l.line) <= 0 {
		return
	}

	var maxHeight float32
	for _, item := range l.line {
		maxHeight = float32(math.Max(float64(maxHeight), float64(fyne.MeasureText(item.Text, item.FontSize, item.FontStyle).Height)))
	}
	baseline := l.cursorY + (maxHeight)*defaultLeading
	for _, item := range l.line {
		l.DisplayList = append(l.DisplayList, DisplayItem{X: item.X, Y: baseline - item.Y, Text: item.Text, FontStyle: item.FontStyle, FontSize: item.FontSize})
	}

	l.line = make([]DisplayItem, 0)
	l.cursorY += maxHeight * defaultLeading
	l.cursorX = DefaultHStep
}

func (l *Layout) createLayout(tokens []interface{}) {
	var inEmoji bool

	for _, tok := range tokens {
		l.token(tok, &inEmoji)
	}

	l.flush()
}
