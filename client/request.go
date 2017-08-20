package client

import (
	"bufio"
	"image"
	"io"
	"io/ioutil"
	//	"log"

	"github.com/as/hub/wire"
	"github.com/as/text"
)

func (c *User) SetOrigin(org int64, exact bool) {
	//log.Printf("SetOrigin 1 %d %b\n", org, exact)
	//log.Printf("SetOrigin 2 %d %b\n", org, exact)
	//log.Printf("SetOrigin 3 %d %b\n", org, exact)

	//log.Printf("SetOrigin %d %b\n", org, exact)
	if !exact {
		//log.Printf("read bytes")
		Data, _ := bufio.NewReader(io.NewSectionReader(c, org, org+512)).ReadBytes('\n')
		//log.Printf("read bytes done")
		org += int64(len(Data))
	}
	//log.Printf("setOrigin %d %b\n", clamp(org, 0, c.Len()), exact)
	c.setOrigin(clamp(org, 0, c.Len()))
	//log.Printf("leave  SetOrigin %d %b\n", org, exact)
}
func (c *User) setOrigin(org int64) {
	//log.Printf("setorigin %d\n", org)
	fl := c.fr.Len()
	switch text.Region5(org, org+fl, c.org, c.org+fl) {
	case -1:
		////c.frInsert(c.Bytes()[org:org+(c.org-org)], 0)
		data, err := ioutil.ReadAll(io.NewSectionReader(c, org, c.org))
		if err != nil {
			//log.Printf("setOrigin: %s\n", err)
		}
		//log.Printf("frInsert: %s\n", data)
		c.frInsert(data, 0)
		//log.Printf("frInsert done: %s\n", data)
		c.org = org
	case -2, 2:
		c.frDelete(0, c.fr.Len())
		c.org = org
		c.Fill()
	case 1:
		c.frDelete(0, org-c.org)
		c.org = org
		c.Fill()
	case 0:
		panic("never happens")
	}
	fr := c.fr.Bounds()
	if pt := c.fr.PointOf(c.fr.Len()); pt.Y != fr.Max.Y {
		c.frPaint(pt, fr.Max, c.fr.Color.Palette.Back)
	}
	Q0, Q1 := c.Dot()
	//log.Printf("c.Select %d:%d\n", Q0, Q1)
	c.Select(Q0, Q1)
}

// Who asks the hub for active user information. It populates
// eaCh users Id, and dot (Q0:Q1), storing them on the
// sIde.
func (c *User) Who() {
	c.ask <- wire.Packet{
		Id:   c.Id,
		Kind: 'w',
	}
	for v := range c.rc {
		ui := &userinfo{Id: v.N, Q0: v.Q0, Q1: v.Q1}
		if ui.Id == -1 {
			break
		}
		//	//log.Printf("userinfo: %#v\n", ui)
		c.knows[v.N] = ui
	}
}

func (c *User) asker() {
	ctr := 0
	for {
		select {
		case ap := <-c.askin:
			ctr++
			ap.RcId = ctr
			c.replyfns[ctr] = replyfn(func(p wire.Packet) {
				ap.rc <- p
			})
			c.ask <- ap.Packet
		case rep := <-c.replyc:
			if fn, ok := c.replyfns[rep.RcId]; ok {
				go fn(rep)
				delete(c.replyfns, rep.RcId)
			} else {
				//log.Printf("no fn for reply #%d\n", rep.RcId)
			}
		}
	}
}

type replyfn func(p wire.Packet)

type request struct {
	rc chan wire.Packet
	wire.Packet
}

func (c *User) writeRequest(p *wire.Packet) (replyc chan wire.Packet) {
	replyc = make(chan wire.Packet)
	p.Id = c.Id
	c.askin <- request{rc: replyc, Packet: *p}
	return replyc
}

func (c *User) Mark(Q0 int64) {
	rc := c.writeRequest(&wire.Packet{
		Kind: 'm',
		Data: wire.Data{
			Q0: Q0,
		},
	})
	//log.Println("waiting for reply on rc")
	<-rc
	return
}
func (c *User) Read(p []byte) (n int, err error) {
	//log.Println("Read: waiting for reply on rc")
	rc := c.writeRequest(&wire.Packet{
		Kind: 'r',
		Data: wire.Data{
			P: p,
		},
	})
	//log.Println("Read: waiting for reply on rc")
	r := <-rc
	copy(p, r.P)
	return r.N, wire.StrToErr(r.Err)
}

func (c *User) ReadAt(p []byte, at int64) (n int, err error) {
	rc := c.writeRequest(&wire.Packet{
		Kind: 'R',
		Data: wire.Data{
			Q0: at,
			P:  p,
		},
	})

	r := <-rc
	copy(p, r.P)
	return r.N, wire.StrToErr(r.Err)
}

func (c *User) Insert(p []byte, Q0 int64) (n int) {
	rc := c.writeRequest(&wire.Packet{
		Id:   c.Id,
		Kind: 'i',
		Data: wire.Data{
			Q0: Q0,
			N:  len(p),
			P:  p,
		},
		RcId: 2,
	})
	//	//log.Println("insert: waiting for reply on rc")
	r := <-rc
	//	//log.Println("insert: got reply")
	return r.N
}
func (c *User) Select(Q0, Q1 int64) {
	rc := c.writeRequest(&wire.Packet{
		Id:   c.Id,
		Kind: 's',
		Data: wire.Data{
			Q0: Q0,
			Q1: Q1,
		},
		RcId: 2,
	})
	//	//log.Println("select: waiting for reply on rc")
	<-rc
	//	//log.Println("select: got reply")
}
func (c *User) Dot() (Q0, Q1 int64) {
	rc := c.writeRequest(&wire.Packet{
		Kind: '.',
	})
	//	//log.Println("dot: waiting for reply on rc")
	r := <-rc
	//	//log.Println("dot: got reply")
	return r.Q0, r.Q1
}
func (c *User) Bytes() (p []byte) {
	rc := c.writeRequest(&wire.Packet{
		Kind: 'b',
		Data: wire.Data{},
	})
	//	//log.Println("byte: waiting for reply on rc")
	r := <-rc
	//	//log.Println("byte: got reply")
	return r.P
}

func (c *User) Delete(Q0, Q1 int64) (n int) {
	rc := c.writeRequest(&wire.Packet{
		Kind: 'd',
		Data: wire.Data{
			Q0: Q0,
			Q1: Q1,
		},
	})
	//	//log.Println("delete: waiting for reply on rc")
	r := <-rc
	//	//log.Println("delete: got reply")
	return r.N
}
func (c *User) Len() (n int64) {
	rc := c.writeRequest(&wire.Packet{
		Kind: 'l',
	})
	//log.Println("len: waiting for reply on rc")
	r := <-rc
	//log.Println("len: got reply")
	return int64(r.N)
}
func (c *User) Fill() {
	/*
		for !c.fr.Full() {
			qep := c.org + c.fr.Len()
			n := min(c.Len()-qep, 2500)
			if n <= 0 {
				break
			}
			rp := c.Bytes()[qep : qep+n]
			nl := c.fr.MaxLine() - c.fr.Line()
			m := 0
			i := int64(0)
			for i < n {
				if rp[i] == '\n' {
					m++
					if m >= nl {
						i++
						break
					}
				}
				i++
			}
			c.frInsert(rp[:i], c.fr.Len())
			//		c.Mark()
		}
	*/
	println("Fill")
	if c.fr == nil {
		return
	}
	for !c.fr.Full() {
		qep := c.Origin() + c.fr.Len()
		n := min(c.Len()-qep, 2500)
		if n <= 0 {
			break
		}
		buf := bufio.NewReader(io.NewSectionReader(c, qep, qep+n))
		for nl := c.fr.MaxLine() - c.fr.Line(); nl >= 0; nl-- {
			line, err := buf.ReadBytes('\n')
			c.frInsert(line, c.fr.Len())
			if err != nil {
				break
			}
		}
	}

}
func (c *User) Origin() int64 {
	println("Origin()")
	return c.org
}
func (c *User) Scroll(dl int) {
	println("Scroll(dl int")
	if dl == 0 {
		return
	}
	org := c.org
	if dl < 0 {
		org = c.BackNL(org, -dl)
		c.SetOrigin(org, true)
	} else {
		if org+c.fr.Len() >= c.Len() {
			return
		}
		r := c.fr.Bounds()
		org += c.fr.IndexOf(image.Pt(r.Min.X, r.Min.Y+dl*c.fr.Font.Dy()))
		c.SetOrigin(org, true)
	}
}

func (c *User) BackNL(p int64, n int) int64 {
	R := c.Bytes()
	if n == 0 && p > 0 && R[p-1] != '\n' {
		n = 1
	}
	for i := n; i > 0 && p > 0; {
		i--
		p--
		if p == 0 {
			break
		}
		for j := 512; j-1 > 0 && p > 0; p-- {
			j--
			if p-1 < 0 || p-1 > c.Len() || R[p-1] == '\n' {
				break
			}
		}
	}
	return p
}
