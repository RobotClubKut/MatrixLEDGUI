package main

import (
	"fmt"
	"log"

	"github.com/mikepb/go-serial"
)

func main() {
	info, err := serial.ListPorts()

	if err != nil {
		log.Fatalln(err)
	}

	for _, p := range info {
		fmt.Println(p.Name())
	}
}
