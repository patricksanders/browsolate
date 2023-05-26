package main

import (
	"fmt"
	"github.com/patricksanders/browsolate"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Print("not enough args")
		os.Exit(2)
	}
	opts := &browsolate.InstanceOpts{}
	err := browsolate.StartIsolatedChromeInstance(os.Args[1], opts)
	if err != nil {
		fmt.Printf("oh no! %v", err)
	}
}
