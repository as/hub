package client

import (
	//"log"

	"github.com/as/frame"
	"github.com/as/hub/wire"
	"github.com/as/text"
)

type userinfo struct {
	Id     int
	Q0, Q1 int64
	col    *frame.Color
}

type User struct {
	Id        int
	Q0, Q1    int64
	askin     chan request
	ask       chan wire.Packet
	replyc    chan wire.Packet
	broadcast chan wire.Packet
	rc        chan wire.Packet
	replyfns  map[int]replyfn
	fr        *frame.Frame
	org       int64
	knows     map[int]*userinfo
	sem       chan bool
	dirty     bool
	text.Sender
}

/*
func (c *User) lock()   { c.sem <- true }
func (c *User) unlock() { <-c.sem }
*/

func (c *User) lock()   {}
func (c *User) unlock() {}

func (c *User) FrameDelete(id int, q0, q1 int64) int {
	return c.frameDelete(q0, q1)
}
func (c *User) FrameInsert(id int, p []byte, q0, q1 int64) int {
	return c.frameInsert(p, q0)
}
func (c *User) FrameSelect(id int, q0, q1 int64) {
	c.frameSelect(id, q0, q1)
}

func (c *User) frameDelete(Q0, Q1 int64) int {
	//log.Printf("in frameDelete\n")
	if Q1 < Q0 {
		return 0
	}
	reg := text.Region5(Q0, Q1, c.org, c.org+c.fr.Len())
	//log.Printf("in region5: %d because Region5(%d,%d,%d,%d)\n", reg, c.org, c.org+c.fr.Len(), Q0, Q1)
	switch reg {
	case -2:
		c.org -= Q1 - Q0
	case -1:
		c.frDelete(0, c.org-Q1)
		c.org = c.BackNL(Q0, 1)
		c.Fill()
	case 0:
		if Q0 < c.org {
			c.frDelete(0, c.fr.Len())
			c.org = c.BackNL(Q0, 1)
			c.Fill()
		} else {
			c.frDelete(Q0-c.org, Q1-c.org)
			//c.Fill()
		}
	case 1:
		c.frDelete(Q0-c.org, c.fr.Len())
		c.Fill()
		//log.Printf("setorg done")
	}
	//log.Printf("ret: %d\n", int(Q1-Q0)+1)
	return int(Q1 - Q0)
}

func (c *User) frameInsert(p []byte, Q0 int64) int {
	return c.frInsert(p, Q0-c.org)
}

func (c *User) frameSelect(id int, q0, q1 int64) {
	var (
		u  *userinfo
		ok bool
	)
	if u, ok = c.knows[id]; !ok {
		u = &userinfo{Id: id, Q0: q0, Q1: q1}
	}
	c.frameSelect0(u, q0, q1, u.Q0, u.Q1)
	u.Q0, u.Q1 = q0, q1
}

func clamp(v, l, h int64) int64 {
	println("clamp")
	if v < l {
		return v
	}
	if v > h {
		return h
	}
	return v
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
