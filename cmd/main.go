package main

import (
	"log"
	"time"
)

func main() {
	for {
		log.Println("data")
		time.Sleep(10 * time.Second)
	}
}
