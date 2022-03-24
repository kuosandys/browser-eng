package browser

import (
	"fmt"
	"image"
	"image/color"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/kuosandys/browser-engineering/pkg/layout"
	"github.com/kuosandys/browser-engineering/pkg/requester"
)

const (
	width      = 800
	height     = 600
	scrollStep = 100
)

type Browser struct {
	window      *g.MasterWindow
	displayList []layout.DisplayItem
	scroll      int
}

// NewBrowser returns a running new browser with some defaults
func NewBrowser() *Browser {
	b := &Browser{
		window: g.NewMasterWindow("hello bello", width, height, g.MasterWindowFlagsNotResizable),
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
	b.window.Run(b.loop)
}

// loop draws the actual content of the browser window
func (b *Browser) loop() {
	g.SingleWindow().RegisterKeyboardShortcuts(
		g.WindowShortcut{Key: g.KeyDown, Callback: b.scrollDown},
	).Layout(
		g.Custom(func() {
			canvas := g.GetCanvas()
			color := color.RGBA{200, 75, 75, 255}
			for _, d := range b.displayList {
				if (d.Y > b.scroll+height) || (d.Y+layout.VStep < b.scroll) {
					continue
				}
				canvas.AddText(image.Pt(d.X, d.Y-b.scroll), color, d.Text)
			}
		}),
	)
}

func (b *Browser) scrollDown() {
	b.scroll += scrollStep
}
