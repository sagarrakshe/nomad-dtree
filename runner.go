package main

import (
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
	NomadClient  *api.Client
	JobsPath     string
	Dependencies []Njob
}

func (r *Runner) get_dependency(currentJob string, body gjson.Result) []Njob {
	var jobs []Njob

	currentJobWait := body.Get(fmt.Sprintf("%s.wait_cond", currentJob))
	preJob := body.Get(fmt.Sprintf("%s.pre.job", currentJob))
	postJob := body.Get(fmt.Sprintf("%s.post.job", currentJob))
	postJobWait := body.Get(fmt.Sprintf("%s.post.wait_cond", currentJob))

	if preJob.Str == "" {
		if postJob.Str == "" {
			return append(jobs, Njob{Job: currentJob, Wait: currentJobWait.Int()})
		}
		return append(jobs, Njob{Job: currentJob, Wait: currentJobWait.Int()},
			Njob{Job: postJob.Str, Wait: postJobWait.Int()})
	}

	t := r.get_dependency(preJob.Str, body)
	t = append(t, Njob{Job: currentJob, Wait: currentJobWait.Int()})
	if postJob.Str != "" {
		return append(t, Njob{Job: postJob.Str, Wait: postJobWait.Int()})
	}
	return t
}

func (r *Runner) run_tree(job string, dependencies []byte) error {
	// Validate Json
	jobs := r.get_dependency(job, gjson.Get(string(dependencies), "dependencies"))
	for idx, j := range jobs {
		log.Printf("Run job: %+v", j.Job)

		skip_wait, err := r.run_job(j.Job)
		if err != nil {
			log.Printf("Error running job: %+v", err)
			return err
		}

		log.Printf("status %+v", skip_wait)
		if !skip_wait {
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

func (r *Runner) run_job(jobName string) (bool, error) {
	nomadJob, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.nomad", r.JobsPath, jobName))
	if err != nil {
		log.Fatalf("error reading file: %+v", err)
		return false, err
	}

	jobs := r.NomadClient.Jobs()
	njob, err := jobs.ParseHCL(string(nomadJob), true)
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
