package main

import (
	"os"

	"github.com/detohm/golox"
)

func main() {

	lox := golox.NewLox()
	lox.Main(os.Args)
}
