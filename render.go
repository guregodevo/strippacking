package strippacking

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
	"os"
)

var draw *gdk.Drawable
var gc *gdk.GC
var drawingarea *gtk.DrawingArea

const W = 6

var H float64 = 5

const MAX_Y = 900
const MAX_X = 1212

func (r *Rect) Draw(filled bool) {
	if 0 == conv_h(r.H) {
		draw.DrawRectangle(gc, filled, conv_x(r.X), conv_y(r.Y+r.H), conv_w(r.W), 1)
		return
	}
	draw.DrawRectangle(gc, filled, conv_x(r.X), conv_y(r.Y+r.H), conv_w(r.W), conv_h(r.H))
}

func conv_y(y float64) int {
	return int(((H - y) / H) * MAX_Y)
}

func conv_x(x float64) int {
	return int((x / W) * MAX_X)
}

func conv_h(h float64) int {
	return int((h / H) * MAX_Y)
}

func conv_w(w float64) int {
	return int((w / W) * MAX_X)
}

func render_all(rects []Rect, m int) {
	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("GTK DrawingArea")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	},
		nil)

	vbox := gtk.NewVBox(true, 0)
	vbox.SetBorderWidth(5)
	drawingarea = gtk.NewDrawingArea()

	var pixmap *gdk.Pixmap

	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		allocation := drawingarea.GetAllocation()
		draw = drawingarea.GetWindow().GetDrawable()
		pixmap = gdk.NewPixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)
		gc = gdk.NewGC(pixmap.GetDrawable())
		gc.SetRgbFgColor(gdk.NewColor("white"))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
		gc.SetRgbFgColor(gdk.NewColor("black"))
		gc.SetRgbBgColor(gdk.NewColor("white"))
	},
		nil)

	drawingarea.Connect("expose-event", func() {
		if pixmap != nil {
			drawingarea.GetWindow().GetDrawable().DrawDrawable(gc,
				pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
		draw_all(rects, m)
	},
		nil)

	vbox.Add(drawingarea)

	window.Add(vbox)
	window.Maximize()
	window.ShowAll()

	gtk.Main()
}

var Pnonsolid *bool
var Prenderbins *bool
var bins_to_render []*Rect = nil

func draw_all(rects []Rect, m int) {
	strip_width := W / float64(m)
	if strip_width > 1 {
		strip_width = 1
	}
	for y := 0; y < m; y++ {
		global := Rect{float64(y) * strip_width, 0, H, strip_width}
		global.Draw(false)
	}
	for _, r := range rects {
		r.W *= strip_width
		r.X *= strip_width
		r.Draw(!*Pnonsolid)
	}
	if *Prenderbins {
		gc.SetRgbFgColor(gdk.NewColor("red"))
		for _, r := range bins_to_render {
			r.W *= strip_width
			r.X *= strip_width
			r.Draw(false)
		}
		gc.SetRgbFgColor(gdk.NewColor("black"))
	}
}