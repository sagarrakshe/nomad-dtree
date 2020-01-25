package main

import (
	"log"
	"os"

	nomad "github.com/hashicorp/nomad/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

const Version = "1.1.0"

var (
	app = kingpin.New("nomad-dtree", "Tool for handling nomad dependencies")

	run = app.Command("run", "Run the nomad jobs")

	stop      = app.Command("stop", "Stop the nomad jobs")
	stopPurge = stop.Flag("purge", "Purge the job").Bool()
	stopDeep  = stop.Flag("deep", "Stop all dependent jobs").Bool()

	job               = app.Flag("job", "Job").String()
	rootFile          = app.Flag("root-ca-file", "RootCA File").Envar("ROOT_CA_FILE").String()
	certFile          = app.Flag("cert-file", "Cert File").Envar("CERT_FILE").String()
	keyFile           = app.Flag("key-file", "Key File").Envar("KEY_FILE").String()
	nomadAddr         = app.Flag("nomad-addr", "Nomad Server Addr").Envar("NOMAD_ADDR").Required().String()
	storeDriver       = app.Flag("store", "store for nomad jobs").Envar("STORE_DRIVER").Default("Filesystem").String()
	consulAddr        = app.Flag("consul-addr", "Consul Address").Envar("CONSUL_ADDRESS").String()
	consulDepFilepath = app.Flag("consul-depfile-path", "Consul Dependency Filepath").Envar("CONSUL_DEP_FILEPATH").String()
	consulJobsPath    = app.Flag("consul-jobs-path", "Consul Jobs path").Envar("CONSUL_JOBS_PATH").String()
	fsDepFilepath     = app.Flag("fs-depfile-path", "Filesystem Dependency File Path").Envar("FS_DEP_FILEPATH").String()
	fsJobsPath        = app.Flag("fs-jobs-path", "Filesystem Path to jobs location").Envar("FS_JOBS_PATH").String()
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	config := &nomad.Config{
		Address: *nomadAddr,
		TLSConfig: &nomad.TLSConfig{
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
		FsDepPath:      *fsDepFilepath,
		FsJobsPath:     *fsJobsPath,
	}

	cmdConfig := &CmdConfig{
		Command:   cmd,
		StopPurge: *stopPurge,
		StopDeep:  *stopDeep,
	}

	runner, err := NewRunner(config, storeConfig, cmdConfig)
	if err != nil {
		log.Fatal(err)
	}

	rerr := runner.run_tree(*job)
	if rerr != nil {
		log.Fatal(rerr)
	}
}
