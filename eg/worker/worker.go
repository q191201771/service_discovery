package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	sd "github.com/q191201771/service_discovery"
)

func main() {
	name := flag.String("name", fmt.Sprintf("%d", time.Now().Unix()), "des")
	extInfo := "nope..."

	flag.Parse()
	w, err := sd.NewWorker("service_d", *name, extInfo, []string{
		"http://127.0.0.1:32379",
		"http://127.0.0.1:22379",
		"http://127.0.0.1:2379",
	})
	if err != nil {
		log.Fatal(err)
	}
	w.Register()
	log.Println("name ->", *name, "extInfo ->", extInfo)
	for {
		log.Println("isActive ->", w.IsActive())
		time.Sleep(time.Second * 2)
	}
}
