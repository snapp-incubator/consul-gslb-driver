package main

import (
	"github.com/sirupsen/logrus"
	"github.com/snapp-cab/consul-gslb-driver/internal/cli"
)

func main() {
	log := logrus.New()
	c := cli.InitCli()

	if err := c.ReadConfig("config.example.yaml"); err != nil {
		log.Fatal(err)
	}

	if err := c.C.Execute(); err != nil {
		log.Fatal(err)
	}
}
