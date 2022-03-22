package main

import (
	"fmt"
	"os"

	"github.com/kuosandys/browser-engineering/pkg/browser"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please input a URL")
		os.Exit(1)
	}

	b := browser.NewBrowser()
	b.Load(args[1])
}
