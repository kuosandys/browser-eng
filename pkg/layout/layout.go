package layout

const (
	hstep = 13
	vstep = 18
)

type DisplayItem struct {
	X    int
	Y    int
	Text string
}

func CreateLayout(text string, width int) []DisplayItem {
	displayList := []DisplayItem{}

	cursorX := hstep
	cursorY := vstep
	for _, c := range text {
		displayList = append(displayList, DisplayItem{X: cursorX, Y: cursorY, Text: string(c)})
		if cursorX >= width-(2*hstep) {
			cursorY += vstep
			cursorX = hstep
		} else {
			cursorX += hstep
		}
	}

	return displayList
}
