package main

import (
	"log"

	"github.com/archsinit/ralph-go/internal/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		log.Fatal(err)
	}
}
