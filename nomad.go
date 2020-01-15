package main

import (
	"log"

	"github.com/hashicorp/nomad/api"
)

func NewNomadClient(config *api.Config) (*api.Client, error) {

	nomadConfig := api.DefaultConfig()
	nomadConfig.Address = config.Address
	nomadConfig.TLSConfig = config.TLSConfig

	cl, err := api.NewClient(nomadConfig)
	if err != nil {
		log.Fatalf("Error creating client %+v", err)
		return nil, err
	}

	return cl, nil
}
