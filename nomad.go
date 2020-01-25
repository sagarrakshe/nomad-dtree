package main

import (
	"log"

	"github.com/hashicorp/nomad/api"
)

const Running = "running"

type NomadRunner struct {
	NomadClient *api.Client
}

func NewNomadRunner(config *api.Config) (*NomadRunner, error) {
	nomadConfig := api.DefaultConfig()
	nomadConfig.Address = config.Address
	nomadConfig.TLSConfig = config.TLSConfig

	cl, err := api.NewClient(nomadConfig)
	if err != nil {
		log.Fatalf("Error creating client %+v", err)
		return nil, err
	}

	return &NomadRunner{NomadClient: cl}, nil
}

func (n *NomadRunner) run(job []byte) (bool, error) {
	jobs := n.NomadClient.Jobs()
	njob, err := jobs.ParseHCL(string(job), true)
	if err != nil {
		log.Printf("error while parsing job hcl: %+v", err)
		return false, err
	}

	info, _, err := jobs.Info(*njob.ID, &api.QueryOptions{})
	if err != nil {
		//NOTE: Handle 404 status code
		log.Printf("Error getting job info: %+v", err)
	} else if *info.Status == Running {
		return true, nil
	}

	resp, _, err := jobs.Register(njob, nil)
	if err != nil {
		log.Fatalf("error registering jobs: %+v", err)
	}
	log.Printf("Success Reponse: %+v", resp)

	return false, nil
}

func (n *NomadRunner) stop(job []byte, purge bool) (bool, error) {
	jobs := n.NomadClient.Jobs()
	njob, err := jobs.ParseHCL(string(job), true)
	if err != nil {
		log.Printf("error while parsing job hcl: %+v", err)
		return false, err
	}

	jId, _, err := jobs.Deregister(*njob.ID, purge, &api.WriteOptions{})
	if err != nil {
		log.Printf("error stopping the job: %+v", err)
		return false, err
	}

	log.Printf("Stopped Job: %+v - %+v", *njob.Name, jId)
	return true, nil
}
