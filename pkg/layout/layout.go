package layout

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

	for _, c := range text {
		char := string(c)

		displayList = append(displayList, DisplayItem{X: cursorX, Y: cursorY, Text: char})

		if char == "\n" {
			cursorY += 2 * VStep
			cursorX = HStep
		} else if cursorX >= width-(2*HStep) {
			cursorY += VStep
			cursorX = HStep
		} else {
			cursorX += HStep
		}
	}

	return displayList
}
