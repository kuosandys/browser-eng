package browser

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/kuosandys/browser-engineering/pkg/layout"
	"github.com/kuosandys/browser-engineering/pkg/requester"
)

const (
	width      = 800
	height     = 600
	scrollStep = 100
)

type Browser struct {
	window         fyne.Window
	displayList    []layout.DisplayItem
	scrollPosition int
}

// NewBrowser returns a running new browser with some defaults
func NewBrowser() *Browser {
	app := app.New()
	window := app.NewWindow("hello bello")
	window.Resize(fyne.NewSize(width, height))

	b := &Browser{
		window: window,
	}
	return b
}

// Load requests the URL and runs the main loop
func (b *Browser) Load(url string) {
	text, err := requester.MakeRequest(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.displayList = layout.CreateLayout(text, width)

	b.window.Canvas().SetOnTypedKey(b.handleKeyEvents)

	b.draw()
	b.window.ShowAndRun()
}

// draw the actual content of the browser window
func (b *Browser) draw() {
	textElements := []fyne.CanvasObject{}

	for _, d := range b.displayList {
		if (d.Y > b.scrollPosition+height) || (d.Y+layout.VStep < b.scrollPosition) {
			continue
		}

		text := canvas.NewText(d.Text, color.White)
		text.Move(fyne.NewPos(float32(d.X), float32(d.Y-b.scrollPosition)))
		textElements = append(textElements, text)
	}

	content := container.NewWithoutLayout(textElements...)
	b.window.SetContent(content)
}

// scroll moves the scroll position
func (b *Browser) scroll(dir int) {
	switch dir {
	case 1:
		if b.scrollPosition == 0 {
			return
		}
		b.scrollPosition -= scrollStep
	case -1:
		b.scrollPosition += scrollStep
	}
}

// handleKeyEvents handles key events
func (b *Browser) handleKeyEvents(keyEvent *fyne.KeyEvent) {
	switch keyEvent.Name {
	case fyne.KeyDown:
		b.scroll(-1)
		b.draw()
	case fyne.KeyUp:
		b.scroll(1)
		b.draw()
	}
}
