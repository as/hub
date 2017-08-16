package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/as/hub"
	"github.com/as/text"
)

func main() {
	a := os.Args
	if len(a) < 2 {
		log.Fatalln("usage: hubs.exe host:port [file]")
	}
	buf := text.NewBuffer()
	var (
		data []byte
		err  error
	)
	if len(a) > 2 {
		data, err = ioutil.ReadFile(a[2])
		if err != nil {
			log.Fatalln(err)
		}
		buf.Insert(data, 0)
	} else {
		buf.Insert([]byte("The quick brown fox"), 0)
	}
	hub := hub.NewHub(buf)
	hub.Run()
	err = hub.ListenAndServe(a[1])
	if err != nil{
		log.Fatalln(err)
	}
	select {}
}
