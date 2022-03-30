package layout

import (
	"strconv"
	"strings"
)

const (
	HStep = 13
	VStep = 18
)

type DisplayItem struct {
	X    int
	Y    int
	Text string
}

func CreateLayout(text string, width int) []DisplayItem {
	displayList := []DisplayItem{}

	cursorX := HStep
	cursorY := VStep
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

		vstep := VStep
		if inEmoji {
			// use two VSteps for emoji unicode characters
			vstep = VStep * 2
		}

		if cursorX >= width-(3*HStep) {
			cursorY += vstep
			cursorX = HStep
		} else {
			cursorX += vstep
		}
	}

	return displayList
}
