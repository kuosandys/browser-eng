package layout

import (
	"strconv"
	"strings"
)

const (
	DefaultHStep float32 = 13
	DefaultVStep float32 = 18
)

type DisplayItem struct {
	X    float32
	Y    float32
	Text string
}

func CreateLayout(text string, width float32, scale float32) []DisplayItem {
	displayList := []DisplayItem{}

	VStep := DefaultVStep * scale
	HStep := DefaultHStep * scale

	cursorX := DefaultHStep
	cursorY := DefaultVStep
	var inEmoji bool

	for _, c := range text {
		char := string(c)

		if char == "\n" {
			cursorY += 2 * VStep
			cursorX = HStep
			continue
		}

		displayList = append(displayList, DisplayItem{X: cursorX, Y: cursorY, Text: char})

		ascii := strings.Trim(strconv.QuoteToASCII(char), "\"")
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

		if cursorX >= width-(3*DefaultHStep) {
			cursorY += VStep
			cursorX = DefaultHStep
		} else {
			cursorX += hStep
		}
	}

	return displayList
}
