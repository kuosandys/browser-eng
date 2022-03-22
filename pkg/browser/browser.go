package browser

import (
	"fmt"
	"image"
	"image/color"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/kuosandys/browser-engineering/pkg/requester"
)

const (
	width  = 800
	height = 600
)

type Browser struct {
	window *g.MasterWindow
	text   string
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

	b.text = text
	b.window.Run(b.loop)
}

// loop draws the actual content of the browser window
func (b *Browser) loop() {
	g.SingleWindow().Layout(
		g.Custom(func() {
			canvas := g.GetCanvas()
			color := color.RGBA{200, 75, 75, 255}
			canvas.AddText(image.Pt(200, 150), color, b.text)
		}),
	)
}
