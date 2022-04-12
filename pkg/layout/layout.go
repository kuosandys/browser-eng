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

type DisplayItem struct {
	X         float32
	Y         float32
	Text      string
	FontStyle fyne.TextStyle
}

func CreateLayout(token []interface{}, width float32, scale float32) []DisplayItem {
	var displayList []DisplayItem

	HStep := DefaultHStep * scale

	cursorX := DefaultHStep
	cursorY := DefaultVStep
	var inEmoji bool

	spaceSize := fyne.MeasureText(" ", HStep, fyne.TextStyle{})
	var fontStyle fyne.TextStyle

	for _, tok := range token {
		switch tok.(type) {
		case *parser.Text:
			for _, word := range strings.Fields(tok.(*parser.Text).Text) {
				displayList = append(displayList, DisplayItem{X: cursorX, Y: cursorY, Text: word, FontStyle: fontStyle})

				ascii := strings.Trim(strconv.QuoteToASCII(word), "\"")
				if len(ascii) > 2 && (ascii[0:2] == "\\U" || ascii[0:2] == "\\u") {
					// don't change cursor position for emoji unicode characters
					inEmoji = true
					continue
				}

				hStep := HStep
				if inEmoji {
					// use two HSteps for emoji unicode characters
					hStep *= 2
				}
				size := fyne.MeasureText(word, HStep, fyne.TextStyle{})

				if cursorX+size.Width >= width-(3*DefaultHStep) {
					cursorY += size.Height * defaultLeading
					cursorX = hStep
				} else {
					cursorX += size.Width + spaceSize.Width
				}
			}
		case *parser.Tag:
			switch tok.(*parser.Tag).Tag {
			case "i":
				fontStyle = fyne.TextStyle{Italic: true}
			case "/i":
				fontStyle = fyne.TextStyle{Italic: false}
			case "b":
				fontStyle = fyne.TextStyle{Bold: true}
			case "/b":
				fontStyle = fyne.TextStyle{Bold: false}
			}
		}

	}

	return displayList
}
