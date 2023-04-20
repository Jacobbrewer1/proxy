package main

import (
	"errors"
	"flag"
	"github.com/jacobbrewer1/reverse-proxy/pkg/filehandler"
	"log"
	"os"
)

func flags() error {
	configLocation := flag.String("config", "", "Location of the config file")

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
		log.Fatalln(err)
	}

	a, err := InitializeApp()
	if err != nil {
		log.Fatalln(err)
	}
	if err := a.start(); err != nil {
		a.logger.Error(err.Error())
		os.Exit(1)
	}
}
