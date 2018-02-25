package client

import (
	"image"

	"github.com/as/frame"
	"github.com/as/text"
	"golang.org/x/exp/shiny/screen"
)

func (c *User) Dirty() bool { return c.dirty || c.fr.Dirty() }

func (c *User) Upload(wind screen.Window, b screen.Buffer, sp image.Point) {
	wind.Upload(sp, b, b.Bounds())
	c.dirty = false
}

func (c *User) frInsert(p []byte, Q0 int64) int {
	if len(p) == 0 || Q0 < 0 || Q0 > c.fr.Len() {
		return 0
	}
	c.lock()
	defer c.unlock()
	return c.fr.Insert(p, Q0)
}
func (c *User) frDelete(p0, p1 int64) int {
	c.lock()
	defer c.unlock()
	return c.fr.Delete(p0, p1)
}
func (c *User) frPaint(p0, p1 image.Point, col image.Image) {
	c.lock()
	defer c.unlock()
	c.fr.Paint(p0, p1, col)
	return
}

func (h *User) frameSelect0(u *userinfo, r0, r1, Q0, Q1 int64) {
	if u.col == nil {
		u.col = &frame.Acme
	}
	h.lock()
	reg := text.Region5(r0, r1, Q0, Q1)
	switch reg {
	case -2, 2:
		h.erase(u.col, Q0, Q1)
		h.adjustIn(u, Q0, Q1)
	case -1:
		h.erase(u.col, r1, Q1)
		h.adjustIn(u, r1, Q1)
	case 1:
		h.erase(u.col, Q0, r0)
		h.adjustIn(u, Q0, r0)
	case 0:
		if !(r0 < Q0) { // in
			h.erase(u.col, Q0, r0)
			h.erase(u.col, r1, Q1)
			h.adjustIn(u, Q0, r0)
			h.adjustIn(u, r1, Q1)
		}
	}
	h.paint(u.col, r0, r1)
	h.unlock()
}

func (h *User) adjustIn(u *userinfo, min, max int64) {

	for _, v := range h.knows {
		if u.Id == v.Id {
			continue
		}
		switch text.Region5(v.Q0, v.Q1, min, max) {
		case -1:
			h.paint(v.col, min, v.Q1)
		case 0:
			h.paint(v.col, v.Q0, v.Q1)
		case 1:
			h.paint(v.col, v.Q0, max)
		}
	}
}

func (c *User) erase(col *frame.Color, Q0, Q1 int64) {
	t := c.fr
	//	log.Printf("hub.erase\n")
	org := c.org
	t.Recolor(t.PointOf(Q0-org), Q0-org, Q1-org, frame.Acme.Palette)
}
func (c *User) paint(col *frame.Color, Q0, Q1 int64) {
	t := c.fr
	//	log.Printf("hub.paint\n")
	org := c.org
	t.Recolor(t.PointOf(Q0-org), Q0-org, Q1-org, col.Hi)
}

/*
	if u.overlapped{
		fmt.Printf("overlapped")
		h.paint(t, u.col, r0, r1)
	} else {
		switch reg{
		case -2, 2:
			h.paint(t, u.col, r0, r1)
		case -1:
			h.paint(t, u.col, r0, Q0)
		case  1:
			h.paint(t, u.col, Q1, r1)
		case 0:
			if r1 > Q1{ // over
			fmt.Printf("over")
				h.paint(t, u.col, r0, Q0)
				h.paint(t, u.col, Q1, r1)
			} else {
				h.paint(t, u.col, r0, r1)
			}
		}
	}
	h.mark(u, r0, r1)
*/
