package main

import (
	"log"
	"os"
)

var opts struct {
	Source string   `long:"source" default:":2203" description:"Source port to listen on"`
	Target []string `long:"target" description:"Target address to forward to"`
	Quiet  bool     `long:"quiet" description:"whether to print logging info or not"`
	Buffer int      `long:"buffer" default:"10240" description:"max buffer size for the socket io"`
}

func main() {
	a, err := InitializeApp()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	a.start()
}
