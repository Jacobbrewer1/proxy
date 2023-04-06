package main

import (
	"errors"
	"flag"
	"github.com/jacobbrewer1/reverse-proxy/pkg/filehandler"
	"log"
	"os"
)

func flags() error {
	configLocation := flag.String("config", "", "The location of the config file")

	flag.Parse()

	if *configLocation == "" {
		return errors.New("no config location provided")
	} else {
		filehandler.Location = *configLocation
	}
	return nil
}

func main() {
	if err := flags(); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	a, err := InitializeApp()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	a.start()
}
