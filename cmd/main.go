package main

import (
	"flag"
	"log"

	"github.com/cutlery47/music-storage/internal/app"
)

func main() {
	test := flag.Bool("test", false, "runs the system tests insted of running a server (if set to true)")
	flag.Parse()

	if *test {
		log.Fatal(app.RunAgent())
	} else {
		log.Fatal(app.Run())
	}
}
