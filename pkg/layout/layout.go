package layout

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"

	"github.com/kuosandys/browser-engineering/pkg/parser"
)

const (
	DefaultHStep   float32 = 13
	DefaultVStep   float32 = 18
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
}

type DisplayItem struct {
	X         float32
	Y         float32
	Text      string
	FontStyle fyne.TextStyle
}

func NewLayout(tokens []interface{}, width float32, scale float32) *Layout {
	l := &Layout{
		DisplayList: make([]DisplayItem, 0),
		cursorX:     DefaultHStep,
		cursorY:     DefaultVStep,
		fontStyle:   fyne.TextStyle{},
		fontSize:    0,
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
		}
	}
}

func (l *Layout) text(token *parser.Text, inEmoji *bool) {
	spaceSize := fyne.MeasureText(" ", DefaultHStep*l.scale, l.fontStyle)

	for _, word := range strings.Fields(token.Text) {
		l.DisplayList = append(l.DisplayList, DisplayItem{X: l.cursorX, Y: l.cursorY, Text: word, FontStyle: l.fontStyle})

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
		size := fyne.MeasureText(word, DefaultHStep*l.scale, l.fontStyle)

		if l.cursorX+size.Width+spaceSize.Width >= l.width-(5*DefaultHStep) {
			l.cursorY += size.Height * defaultLeading
			l.cursorX = hStep
		} else {
			l.cursorX += size.Width + spaceSize.Width
		}
	}
}

func (l *Layout) createLayout(tokens []interface{}) {
	var inEmoji bool

	for _, tok := range tokens {
		l.token(tok, &inEmoji)
	}
}
