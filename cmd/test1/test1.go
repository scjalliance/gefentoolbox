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
	for {
		s.RawCommand("r 2 1")
		time.Sleep(time.Second * 10)
		s.RawCommand("r 1 1")
		time.Sleep(time.Second * 10)
	}
}
