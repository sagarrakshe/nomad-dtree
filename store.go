package main

type Store interface {
	GetJob(string) ([]byte, error)
	GetDependencies() ([]byte, error)
}

type StoreConfig struct {
	Driver         string
	ConsulAddr     string
	ConsulJobsPath string
	ConsulDepPath  string
	FsJobsPath     string
	FsDepPath      string
}

func NewStoreClient(config *StoreConfig) (Store, error) {
	if config.Driver == "filesystem" {
		return NewFileSystemClient(config)
	} else {
		return NewConsulClient(config)
	}
}
