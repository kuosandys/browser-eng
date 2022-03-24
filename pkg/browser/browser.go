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
	width  = 800
	height = 600
)

type Browser struct {
	window      *g.MasterWindow
	displayList []layout.DisplayItem
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
	b.window.Run(b.draw)
}

// draw draws the actual content of the browser window
func (b *Browser) draw() {
	g.SingleWindow().Layout(
		g.Custom(func() {
			canvas := g.GetCanvas()
			color := color.RGBA{200, 75, 75, 255}
			for _, d := range b.displayList {
				canvas.AddText(image.Pt(d.X, d.Y), color, d.Text)
			}

		}),
	)
}
