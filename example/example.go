package main

import (
	"fmt"
	"github.com/as/hub"
	"github.com/as/text"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	a := os.Args
	if len(a) < 2 {
		log.Fatalln("usage: hubs.exe host:port")
	}
	buf := text.NewBuffer()
	var (
		data []byte
		err  error
	)
	for _, name := range []string{`\windows\system32\drivers\etc\hosts`, `/etc/hosts`, `/ndb/local`} {
		data, err = ioutil.ReadFile(name)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Fatalln(fmt.Errorf("no hosts files\n"))
	}
	buf.Insert(data, 0)
	hub := hub.NewHub(buf)
	hub.Run()
	hub.ListenAndServe(a[1])
	select {}
}
