package main

import (
	"log"
	"os"

	"github.com/detohm/golox/astgen"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("invalid argument")
		return
	}
	astgen.Generate(os.Args[1])
}
