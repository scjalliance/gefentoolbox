package main

import (
	"fmt"
	"log"
	"time"

	"github.com/scjalliance/gefentoolbox"
)

var addr = "10.100.0.60:23"
var username = "Admin"
var password = "123456"

func main() {
	s, err := gefentoolbox.New(addr, username, password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s.Welcome())
	fmt.Printf("Model: %s\n", s.Model())
	for {
		for in := 1; in < 5; in++ {
			fmt.Println("-----")
			raw, err := s.Route(in, 1)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Printf("Raw result: %s\n", raw)
			}
			time.Sleep(time.Second)
			for i := 1; ; i++ {
				r, err := s.GetRoute(i)
				if err != nil {
					break
				} else {
					fmt.Printf("Route: %d -> %d (%s -> %s)\n", r.Input, r.Output, r.InputName, r.OutputName)
				}
			}
			time.Sleep(time.Second * 5)
		}
	}
}
