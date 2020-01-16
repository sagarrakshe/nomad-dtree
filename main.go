package main

import (
	"log"

	"github.com/hashicorp/nomad/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

const Version = "1.0.0"

var (
	job            = kingpin.Flag("job", "Job").Short('j').String()
	rootFile       = kingpin.Flag("root-ca-file", "RootCA File").Envar("ROOT_CA_FILE").String()
	certFile       = kingpin.Flag("cert-file", "Cert File").Envar("CERT_FILE").String()
	keyFile        = kingpin.Flag("key-file", "Key File").Envar("KEY_FILE").String()
	dependencyFile = kingpin.Flag("dependency-file", "Dependency File").Envar("DEPENDENCY_FILE").String()
	jobsPath       = kingpin.Flag("jobs-path", "Path to jobs location").Envar("JOBS_PATH").String()
	serverAddr     = kingpin.Flag("server-addr", "Server Addr").Short('s').Default("http://127.0.0.1:4646").Envar("SERVER_ADDR").String()
)

func main() {
	kingpin.Parse()

	config := &api.Config{
		Address: *serverAddr,
		TLSConfig: &api.TLSConfig{
			CACert:     *rootFile,
			ClientCert: *certFile,
			ClientKey:  *keyFile,
			Insecure:   true,
		},
	}

	runner, err := NewRunner(config, *dependencyFile, *jobsPath)
	if err != nil {
		log.Fatal(err)
	}

	rerr := runner.run_tree(*job)
	if rerr != nil {
		log.Fatal(rerr)
	}
}
