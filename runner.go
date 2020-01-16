package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/tidwall/gjson"
)

type Njob struct {
	Job  string `json:"job"`
	Wait int64  `json:"wait"`
}

type Runner struct {
	NomadRunner        *NomadRunner
	JobsPath           string
	Dependencies       []Njob
	DependencyFilePath string
}

func NewRunner(nomadConfig *api.Config, dependencyFile string, jobsPath string) (*Runner, error) {
	// Create Nomad client
	nomadRunner, err := NewNomadRunner(nomadConfig)
	if err != nil {
		log.Fatalf("error creating nomad runner: %+v", err)
		return nil, err
	}

	return &Runner{
		NomadRunner:        nomadRunner,
		JobsPath:           jobsPath,
		DependencyFilePath: dependencyFile,
	}, nil
}

func (r *Runner) _dependency(currentJob string, body gjson.Result, idx int) ([]Njob, error) {
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
			return append(jobs, Njob{Job: currentJob, Wait: currentJobWait.Int()}), nil
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
	dependencies, err := ioutil.ReadFile(r.DependencyFilePath)
	if err != nil {
		log.Fatalf("error reading file: %+v", err)
		return nil, err
	}

	return r._dependency(currentJob,
		gjson.Get(string(dependencies), "dependencies"), 5)
}

func (r *Runner) run_tree(job string) error {
	// Get all dependencies of the job
	jobs, err := r.get_dependency(job)
	if err != nil {
		log.Fatalf("error generarting dependency: %+v", err)
		return err
	}

	for idx, j := range jobs {
		log.Printf("Run job: %+v", j.Job)

		nomadJob, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.nomad", r.JobsPath, j.Job))
		if err != nil {
			log.Fatalf("error reading file: %+v", err)
			return err
		}

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
	}

	return nil
}
