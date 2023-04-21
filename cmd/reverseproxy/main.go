package main

import (
	"log"
	"os"
)

func main() {
	a, err := InitializeApp()
	if err != nil {
		log.Fatalln(err)
	}
	if err := a.start(); err != nil {
		a.logger.Error(err.Error())
		os.Exit(1)
	}
}
