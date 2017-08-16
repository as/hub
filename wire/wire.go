package wire

import (
	"fmt"

	"github.com/as/text"
)

const (
	Broadcast = 10
	Reply     = 11
)

type Packet struct {
	Id   int
	Kind byte
	RcId int
	Data
}

type Note struct {
	Ch int
	Packet
}

type Data struct {
	Q0, Q1 int64
	N      int
	Err    string
	P      []byte
}

func StrToErr(str string) error {
	if str == "" {
		return nil
	}
	return fmt.Errorf("%s", str)
}
func Err(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func (p Packet) String() string {
	return fmt.Sprintf("packet %d: kind: '%c', rc: %d, dot: (%d:%d), n: %d, len(p): %d, err: %s\n",
		p.Id, p.Kind, p.RcId, p.Q0, p.Q1, p.N, len(p.P), p.Err)
}

func PacketOk(Kind byte) bool {
	return true
	var valIdKind = []byte("IdswrR.blmELOCS")
	return text.Any(Kind, valIdKind) != -1
}
