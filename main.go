package main

import (
	"log"
	"os"

	"github.com/hashicorp/nomad/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

const Version = "1.0.0"

var (
	app        = kingpin.New("nomad-dtree", "Tool for handling nomad dependencies")
	job        = app.Flag("job", "Job").Short('j').String()
	rootFile   = app.Flag("root-ca-file", "RootCA File").Envar("ROOT_CA_FILE").String()
	certFile   = app.Flag("cert-file", "Cert File").Envar("CERT_FILE").String()
	keyFile    = app.Flag("key-file", "Key File").Envar("KEY_FILE").String()
	serverAddr = app.Flag("server-addr", "Server Addr").Short('s').Default("http://127.0.0.1:4646").Envar("SERVER_ADDR").String()

	storeDriver       = app.Flag("store", "store for nomad jobs").Envar("STORE_DRIVER").Default("Filesystem").String()
	consulAddr        = app.Flag("consul-addr", "Consul Address").Envar("CONSUL_ADDRESS").String()
	consulDepFilepath = app.Flag("consul-depfile-path", "Consul Dependency Filepath").Envar("CONSUL_DEP_FILEPATH").String()
	consulJobsPath    = app.Flag("consul-jobs-path", "Consul Jobs path").Envar("CONSUL_JOBS_PATH").String()
	fsDepFilepath     = app.Flag("fs-depfile-path", "Filesystem Dependency File Path").Envar("FS_DEP_FILEPATH").String()
	fsJobsPath        = app.Flag("fs-jobs-path", "Filesystem Path to jobs location").Envar("FS_JOBS_PATH").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	config := &api.Config{
		Address: *serverAddr,
		TLSConfig: &api.TLSConfig{
			CACert:     *rootFile,
			ClientCert: *certFile,
			ClientKey:  *keyFile,
			Insecure:   true,
		},
	}

	storeConfig := &StoreConfig{
		Driver:         *storeDriver,
		ConsulAddr:     *consulAddr,
		ConsulDepPath:  *consulDepFilepath,
		ConsulJobsPath: *consulJobsPath,
		FsJobsPath:     *fsDepFilepath,
		FsDepPath:      *fsJobsPath,
	}

	runner, err := NewRunner(config, storeConfig)
	if err != nil {
		log.Fatal(err)
	}

	rerr := runner.run_tree(*job)
	if rerr != nil {
		log.Fatal(rerr)
	}
}
