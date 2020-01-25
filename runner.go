package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/tidwall/gjson"
)

const (
	Run  = "run"
	Stop = "stop"

	MaxRecursionDepth = 100
)

type Njob struct {
	Job  string `json:"job"`
	Wait int64  `json:"wait"`
}

type CmdConfig struct {
	Command   string
	StopPurge bool
	StopDeep  bool
}

type Runner struct {
	NomadRunner  *NomadRunner
	Dependencies []Njob
	StoreClient  Store
	Cmd          *CmdConfig
}

func NewRunner(nomadConfig *api.Config, storeConfig *StoreConfig,
	cmdConfig *CmdConfig) (*Runner, error) {

	// Create Nomad client
	nomadRunner, err := NewNomadRunner(nomadConfig)
	if err != nil {
		log.Fatalf("error creating nomad runner: %+v", err)
		return nil, err
	}

	storeClient, err := NewStoreClient(storeConfig)
	if err != nil {
		log.Fatalf("error creating store client: %+v", err)
		return nil, err
	}

	return &Runner{
		NomadRunner: nomadRunner,
		StoreClient: storeClient,
		Cmd:         cmdConfig,
	}, nil
}

func (r *Runner) _dependency(currentJob string, body gjson.Result,
	idx int) ([]Njob, error) {

	var jobs []Njob

	// NOTE: Infinite recursion
	if idx == 0 {
		return nil, errors.New("Infinite recursion")
	}

	// Missing dependency
	if len(body.Get(currentJob).Map()) == 0 {
		return nil, errors.New(fmt.Sprintf("Missing dependency: %s", currentJob))
	}

	currentJobWait := body.Get(fmt.Sprintf("%s.wait_cond", currentJob))
	preJob := body.Get(fmt.Sprintf("%s.pre.job", currentJob))
	postJob := body.Get(fmt.Sprintf("%s.post.job", currentJob))
	postJobWait := body.Get(fmt.Sprintf("%s.post.wait_cond", currentJob))

	if preJob.Str == "" {
		if postJob.Str == "" {
			return append(jobs, Njob{Job: currentJob,
				Wait: currentJobWait.Int()}), nil
		}
		return append(jobs, Njob{Job: currentJob, Wait: currentJobWait.Int()},
			Njob{Job: postJob.Str, Wait: postJobWait.Int()}), nil
	}

	t, err := r._dependency(preJob.Str, body, (idx - 1))
	if err != nil {
		return nil, err
	}
	t = append(t, Njob{Job: currentJob, Wait: currentJobWait.Int()})
	if postJob.Str != "" {
		return append(t, Njob{Job: postJob.Str, Wait: postJobWait.Int()}), nil
	}

	return t, nil
}

func (r *Runner) get_dependency(currentJob string) ([]Njob, error) {
	// Read the dependency JSON file
	dep, err := r.StoreClient.GetDependencies()
	if err != nil {
		return nil, err
	}

	return r._dependency(currentJob,
		gjson.Get(string(dep), "dependencies"), MaxRecursionDepth)
}

func (r *Runner) run_tree(job string) error {
	// Get all dependencies of the job
	jobs, err := r.get_dependency(job)
	if err != nil {
		log.Fatalf("error generarting dependency: %+v", err)
		return err
	}

	if r.Cmd.Command == Stop && !r.Cmd.StopDeep {
		jobs = jobs[len(jobs)-1:]
	}

	for idx, j := range jobs {
		nomadJob, err := r.StoreClient.GetJob(j.Job)
		if err != nil {
			log.Fatalf("error reading file: %+v", err)
			return err
		}

		if r.Cmd.Command == Run {
			skip_wait, err := r.NomadRunner.run(nomadJob)
			if err != nil {
				log.Printf("Error running job: %+v", err)
				return err
			}

			log.Printf("Job Status %+v", skip_wait)

			if !skip_wait {
				// If last job is submitted no need to wait
				if !(idx == len(jobs)-1) {
					log.Printf("Wait for %+v seconds", j.Wait)
					time.Sleep(time.Duration(j.Wait) * time.Second)
				}
			} else {
				log.Printf("Skip Wait, job already running.")
			}
		} else if r.Cmd.Command == Stop {
			_, err := r.NomadRunner.stop(nomadJob, r.Cmd.StopPurge)
			if err != nil {
				log.Printf("Error stop job: %+v", err)
				return err
			}
		}
	}

	return nil
}
