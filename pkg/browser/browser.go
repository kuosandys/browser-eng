package browser

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/kuosandys/browser-engineering/pkg/layout"
	"github.com/kuosandys/browser-engineering/pkg/requester"
)

const (
	width      float32 = 800
	height     float32 = 600
	scrollStep float32 = 100
)

type Browser struct {
	window         fyne.Window
	text           []interface{}
	displayList    []layout.DisplayItem
	scrollPosition float32
	scale          float32
}

// NewBrowser returns a running new browser with some defaults
func NewBrowser() *Browser {
	a := app.New()
	window := a.NewWindow("hello bello")
	window.Resize(fyne.NewSize(width, height))

	b := &Browser{
		window: window,
	}

	b.scale = 1

	return b
}

// Load requests the URL and runs the main loop
func (b *Browser) Load(url string) {
	text, err := requester.MakeRequest(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.text = text
	b.displayList = layout.NewLayout(text, width, b.scale).DisplayList

	b.window.Canvas().SetOnTypedKey(b.handleKeyEvents)

	b.draw()
	b.window.ShowAndRun()
}

// draw the actual content of the browser window
func (b *Browser) draw() {
	elements := b.parseDisplayListToCanvasElements()
	content := container.NewWithoutLayout(elements...)
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

func (b *Browser) magnify(in int) {
	switch in {
	case 1:
		b.scale += 0.1
	case -1:
		b.scale -= 0.1
	}
	b.displayList = layout.NewLayout(b.text, width, b.scale).DisplayList
}

// handleKeyEvents handles key events
func (b *Browser) handleKeyEvents(keyEvent *fyne.KeyEvent) {
	switch keyEvent.Name {
	case fyne.KeyDown:
		b.scroll(-1)
	case fyne.KeyUp:
		b.scroll(1)
	case fyne.KeyEqual:
		b.magnify(1)
	case fyne.KeyMinus:
		b.magnify(-1)
	}
	b.draw()
}

func (b *Browser) parseDisplayListToCanvasElements() []fyne.CanvasObject {
	var elements []fyne.CanvasObject

	for _, d := range b.displayList {
		if (d.Y > b.scrollPosition+height) || (d.Y+layout.DefaultVStep < b.scrollPosition) {
			continue
		}

		ascii := strings.ToUpper(strings.Trim(strconv.QuoteToASCII(d.Text), "\""))

		if strings.Contains(ascii, "\\U") {
			// handle emoji
			startIndices := make([]int, 0)
			i := 0
			for {
				searchStr := ascii[i:]
				startIdx := strings.Index(searchStr, "\\U")
				if startIdx == -1 {
					break
				}
				startIndices = append(startIndices, startIdx+i)
				i += startIdx + 2
			}
			var filePathParts []string
			for i := range startIndices {
				var endIdx int
				if i == (len(startIndices) - 1) {
					endIdx = len(ascii)
				} else {
					endIdx = startIndices[i+1]
				}

				reg, err := regexp.Compile("[^a-zA-Z0-9]+")
				if err != nil {
					log.Fatal(err)
				}
				filePathPart := reg.ReplaceAllString(strings.TrimPrefix(strings.TrimPrefix(ascii[startIndices[i]:endIdx], "\\U"), "000"), "")
				filePathParts = append(filePathParts, filePathPart)

			}

			var filePaths []string
			filePath := os.Getenv("EMOJI_DIR") + strings.Join(filePathParts, "-") + ".png"
			if _, err := os.Stat(filePath); err == nil {
				filePaths = append(filePaths, filePath)
			} else {
				// if the file at the complete path does not exist, try to parse each part as a path
				for _, filePathPart := range filePathParts {
					filePath := os.Getenv("EMOJI_DIR") + filePathPart + ".png"
					if _, err := os.Stat(filePath); err == nil {
						filePaths = append(filePaths, filePath)
					}
				}
			}

			for _, filePath := range filePaths {
				img := canvas.NewImageFromFile(filePath)
				img.Move(fyne.NewPos(d.X-layout.DefaultHStep/2, d.Y-layout.DefaultVStep/4-b.scrollPosition))
				img.Resize(fyne.NewSize(layout.DefaultHStep*2*b.scale, layout.DefaultVStep*2*b.scale))
				img.FillMode = canvas.ImageFillContain
				elements = append(elements, img)
			}
			continue
		}

		// handle text
		text := canvas.NewText(d.Text, color.White)
		text.TextSize = d.Font.Size
		text.TextStyle = d.Font.Style
		text.Move(fyne.NewPos(d.X, d.Y-b.scrollPosition))
		elements = append(elements, text)

	}

	return elements
}
