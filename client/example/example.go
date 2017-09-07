package main

import (
	"github.com/as/event"
	"github.com/as/frame"
	"github.com/as/frame/font"
	"github.com/as/hub/client"
	kbd "github.com/as/text/kbd"
	mous "github.com/as/text/mouse"
	"github.com/as/text/win"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"os"
	"strconv"
	"github.com/as/text"
	"time"
)

var (
	winSize = image.Pt(1800, 1000)
	tagY    = 16
	fontdy  = 16
	focused = false
	redraw  = false
)

func p(e mouse.Event) image.Point {
	return image.Pt(int(e.X), int(e.Y))
}

func atoi(a string) int64 {
	x, _ := strconv.Atoi(a)
	return int64(x)
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("usage: example host:port")
	}
	driver.Main(func(src screen.Screen) {
		wind, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y, "Win"})
		pad := image.Pt(fontdy, fontdy)
		b, err := src.NewBuffer(winSize)
		if err != nil {
			log.Fatalln(err)
		}
		wind.Upload(image.ZP, b, b.Bounds())
		sp := image.ZP
		ed,_ := text.Open(text.NewBuffer())
		w := win.New(sp, pad, b.RGBA(),ed, font.NewGoMono(fontdy))
		wind.Upload(sp, b, b.Bounds())
		wind.Send(paint.Event{})

		mousein := mous.NewMouse(time.Second/3, wind)

		col := frame.Acme
		col.Hi.Back = image.NewUniform(color.RGBA{0x99, 0xDD, 0x99, 192})
		c2 := client.DialEvent(rand.Int(), w.Frame, wind, "tcp", os.Args[1])
		if c2 == nil {
			log.Fatalln("cant connect")
		}
		c2.SetOrigin(0, true)
		var (
			q0, q1, s int64
			r         = w.Bounds()
		)
		mousein.Machine.SetRect(image.Rect(r.Min.X, r.Min.Y+pad.Y, r.Max.X, r.Max.Y-pad.Y))
		redraw := false
		ckRedraw := func() {
			if c2.Dirty() || redraw {
				wind.Send(paint.Event{})
				redraw = false
			}
		}
		for {
			switch e := wind.NextEvent().(type) {
			case mouse.Event:
				e.X -= float32(sp.X)
				e.Y -= float32(sp.Y)
				mousein.Sink <- e
			case mous.Drain:
			DrainLoop:
				for {
					switch wind.NextEvent().(type) {
					case mous.DrainStop:
						break DrainLoop
					}
				}
			case event.Event:
				switch e := e.(type) {
				case event.Insert:
					c2.FrameInsert(e.ID, e.P, e.Q0, e.Q1)
				case event.Delete:
					c2.FrameDelete(e.ID, e.Q0, e.Q1)
				case event.Select:
					c2.FrameSelect(e.ID, e.Q0, e.Q1)
				}
				ckRedraw()
			case mous.SweepEvent:
				s, q0, q1 = mous.Sweep(w, e, pad.Y, s, q0, q1, wind)
				if e.Button == 1 {
					//w.Select(q0, q1)
				}
				ckRedraw()
			case mous.MarkEvent:
				q0 = c2.Origin() + w.IndexOf(p(e.Event))
				q1 = q0
				s = q0
				if e.Button == 1 {
					c2.Select(q0, q1)
					redraw = true
				}
				ckRedraw()
			case mous.ClickEvent:
				switch e.Button {
				case 1:
					c2.Select(q0, q1)
					redraw = true
				}
				ckRedraw()
			case mous.SelectEvent:
				if e.Button == 1 {
					c2.Select(q0, q1)
					redraw = true
				}
				ckRedraw()
			case key.Event:
				if e.Direction == 2 {
					continue
				}
				kbd.SendClient(c2, e)
				ckRedraw()
			case paint.Event:
				c2.Upload(wind, b, sp)
				wind.Publish()
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
				// NT doesn't repaint the window if another window covers it
				if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff {
					focused = false
				} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
					focused = true
				}
			}
		}
	})
}

func region(c, p0, p1 int64) int {
	if c < p0 {
		return -1
	}
	if c >= p1 {
		return 1
	}
	return 0
}

func drawBorder(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, thick int) {
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+thick), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-thick, r.Max.X, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+thick, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Max.X-thick, r.Min.Y, r.Max.X, r.Max.Y), src, sp, draw.Src)
}
