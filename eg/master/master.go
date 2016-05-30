package main

import (
	"log"
	"time"

	sd "github.com/q191201771/service_discovery"
)

func main() {
	m, err := sd.NewMaster("service_d", []string{
		"http://127.0.0.1:2379",
		"http://127.0.0.1:22379",
		"http://127.0.0.1:32379",
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		log.Println("all ->", m.GetNodes())
		log.Println("all(strictly) ->", m.GetNodesStrictly())
		time.Sleep(time.Second * 2)
	}
}
