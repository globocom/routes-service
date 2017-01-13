package main

import (
	"log"

	"github.com/pkg/errors"
)

func main() {
	bgpService, err := NewBgpService(Config{
		EtcdEndpoints: []string{"192.168.33.10:2379"},
		GobgpAddress:  "192.168.33.10:50051",
	})
	if err != nil {
		log.Fatal(err)
	}

	api, err := NewAPI(bgpService)
	if err != nil {
		log.Fatal(errors.Wrap(err, "initialize API failed"))
	}

	bgpService.Start()

	// Blocking
	api.Start()
}
