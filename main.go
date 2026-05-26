package main

import (
	"log"

	"github.com/archsinit/ralph-go/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
