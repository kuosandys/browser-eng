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
	url    string
}

// NewBrowser returns a running new browser with some defaults
func NewBrowser() *Browser {
	b := &Browser{
		window: g.NewMasterWindow("hello bello", width, height, g.MasterWindowFlagsNotResizable),
	}
	return b
}

// Load sets the url of the browser and runs the main loop
func (b *Browser) Load(url string) {
	b.url = url
	b.window.Run(b.loop)
}

// loop draws the actual content of the browser window
func (b *Browser) loop() {
	g.SingleWindow().Layout(
		g.Custom(func() {
			canvas := g.GetCanvas()
			pos := g.GetCursorScreenPos()
			color := color.RGBA{200, 75, 75, 255}
			canvas.AddCircleFilled(pos.Add(image.Pt(100, 150)), 50, color)
			canvas.AddText(image.Pt(200, 150), color, b.url)
		}),
	)
}
