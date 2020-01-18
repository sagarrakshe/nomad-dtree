package main

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/magiconair/properties/assert"
)

func TestRunner(t *testing.T) {
	t.Run("Test Filesystem", func(t *testing.T) {
		sc := &StoreConfig{
			Driver:     "filesystem",
			FsJobsPath: "testdata/",
		}

		t.Run("Test job dependency order", func(t *testing.T) {
			sc.FsDepPath = "testdata/test.json"
			runner, err := NewRunner(api.DefaultConfig(), sc)
			if err != nil {
				t.Fatal(err)
			}

			expected := []Njob{}
			expected = append(expected, Njob{Job: "postgres", Wait: 1})
			expected = append(expected, Njob{Job: "api", Wait: 2})

			actual, err := runner.get_dependency("api")
			assert.Equal(t, expected, actual)
		})

		t.Run("Test cyclic job dependency", func(t *testing.T) {
			sc.FsDepPath = "testdata/recursion.json"
			runner, err := NewRunner(api.DefaultConfig(), sc)
			if err != nil {
				t.Fatal(err)
			}

			expected := "Infinite recursion"
			_, merr := runner.get_dependency("nginx")
			assert.Equal(t, expected, merr.Error())
		})

		t.Run("Test fail when dependent job doesn't exists", func(t *testing.T) {
			sc.FsDepPath = "testdata/missing_dependency.json"
			runner, err := NewRunner(api.DefaultConfig(), sc)
			if err != nil {
				t.Fatal(err)
			}

			expected := "Missing dependency: backend"
			_, merr := runner.get_dependency("nginx")
			assert.Equal(t, expected, merr.Error())
		})
	})
}
