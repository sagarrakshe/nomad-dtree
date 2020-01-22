package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/hashicorp/nomad/api"
	"github.com/magiconair/properties/assert"
)

// Populate the testdata
func populate_consul(sc *StoreConfig) error {
	c, err := NewConsulClient(sc)
	if err != nil {
		return err
	}

	kv := c.Client.KV()
	for _, f := range []string{"test.json", "recursion.json",
		"missing_dependency.json"} {

		file, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", f))
		if err != nil {
			return err
		}
		_, err = kv.Put(&consul.KVPair{Key: f, Value: file},
			&consul.WriteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func TestRunner(t *testing.T) {
	cmd := &CmdConfig{
		Command: "run",
	}
	t.Run("Test Filesystem", func(t *testing.T) {
		sc := &StoreConfig{
			Driver:     "filesystem",
			FsJobsPath: "testdata/",
		}

		t.Run("Test job dependency order", func(t *testing.T) {
			sc.FsDepPath = "testdata/test.json"
			runner, err := NewRunner(api.DefaultConfig(), sc, cmd)
			if err != nil {
				t.Fatal(err)
			}

			expected := []Njob{}
			expected = append(expected, Njob{Job: "postgres", Wait: 1})
			expected = append(expected, Njob{Job: "api", Wait: 2})
			expected = append(expected, Njob{Job: "nginx", Wait: 5})

			actual, err := runner.get_dependency("nginx")
			assert.Equal(t, expected, actual)
		})

		t.Run("Test cyclic job dependency", func(t *testing.T) {
			sc.FsDepPath = "testdata/recursion.json"
			runner, err := NewRunner(api.DefaultConfig(), sc, cmd)
			if err != nil {
				t.Fatal(err)
			}

			expected := "Infinite recursion"
			_, merr := runner.get_dependency("nginx")
			assert.Equal(t, expected, merr.Error())
		})

		t.Run("Test fail when dependent job doesn't exists", func(t *testing.T) {
			sc.FsDepPath = "testdata/missing_dependency.json"
			runner, err := NewRunner(api.DefaultConfig(), sc, cmd)
			if err != nil {
				t.Fatal(err)
			}

			expected := "Missing dependency: backend"
			_, merr := runner.get_dependency("nginx")
			assert.Equal(t, expected, merr.Error())
		})
	})

	t.Run("Test Consul", func(t *testing.T) {
		sc := &StoreConfig{
			Driver:     "consul",
			ConsulAddr: "http://localhost:8500",
		}

		err := populate_consul(sc)
		if err != nil {
			t.Fatal(err)
		}

		t.Run("Test job dependency order", func(t *testing.T) {
			sc.ConsulDepPath = "test.json"
			runner, err := NewRunner(api.DefaultConfig(), sc, cmd)
			if err != nil {
				t.Fatal(err)
			}

			expected := []Njob{}
			expected = append(expected, Njob{Job: "postgres", Wait: 1})
			expected = append(expected, Njob{Job: "api", Wait: 2})
			expected = append(expected, Njob{Job: "nginx", Wait: 5})

			actual, err := runner.get_dependency("nginx")
			assert.Equal(t, expected, actual)
		})

		t.Run("Test cyclic job dependency", func(t *testing.T) {
			sc.ConsulDepPath = "recursion.json"
			runner, err := NewRunner(api.DefaultConfig(), sc, cmd)
			if err != nil {
				t.Fatal(err)
			}

			expected := "Infinite recursion"
			_, merr := runner.get_dependency("nginx")
			assert.Equal(t, expected, merr.Error())
		})

		t.Run("Test fail when dependent job doesn't exists", func(t *testing.T) {
			sc.ConsulDepPath = "missing_dependency.json"
			runner, err := NewRunner(api.DefaultConfig(), sc, cmd)
			if err != nil {
				t.Fatal(err)
			}

			expected := "Missing dependency: backend"
			_, merr := runner.get_dependency("nginx")
			assert.Equal(t, expected, merr.Error())
		})
	})
}
