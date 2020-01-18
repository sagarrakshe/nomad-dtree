package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

type Filesystem struct {
	DepFilepath string
	JobsPath    string
}

func NewFileSystemClient(config *StoreConfig) (*Filesystem, error) {
	if config.FsDepPath == "" {
		return nil, errors.New("Dependency filepath empty")
	}

	if config.FsJobsPath == "" {
		return nil, errors.New("Jobs filepath empty")
	}

	return &Filesystem{DepFilepath: config.FsDepPath,
		JobsPath: config.FsJobsPath}, nil
}

func (f *Filesystem) GetDependencies() ([]byte, error) {
	// NOTE: Add Validate if the file exists
	dependencies, err := ioutil.ReadFile(f.DepFilepath)
	if err != nil {
		log.Fatalf("error reading dependency file: %+v", err)
		return nil, err
	}
	return dependencies, nil
}

func (f *Filesystem) GetJob(job string) ([]byte, error) {
	// NOTE: Add Validate if the file exists
	nomadJob, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.nomad", f.JobsPath, job))
	if err != nil {
		log.Fatalf("error reading file: %+v", err)
		return nil, err
	}
	return nomadJob, nil
}
