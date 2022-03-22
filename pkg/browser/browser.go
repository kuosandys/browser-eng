package browser

import (
	"image"
	"image/color"

	g "github.com/AllenDang/giu"
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

// Load sets the url of the browser and runs the main loop
func (b *Browser) Load(text string) {
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
