package main

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {

	t.Run("Test job dependency order", func(t *testing.T) {
		runner, err := NewRunner(api.DefaultConfig(), "testdata/test.json", "")
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
		runner, err := NewRunner(api.DefaultConfig(), "testdata/recursion.json", "")
		if err != nil {
			t.Fatal(err)
		}

		expected := "Infinite recursion"
		_, merr := runner.get_dependency("nginx")
		assert.Equal(t, expected, merr.Error())
	})

	t.Run("Test fail when dependent job doesn't exists", func(t *testing.T) {
		runner, err := NewRunner(api.DefaultConfig(),
			"testdata/missing_dependency.json", "")
		if err != nil {
			t.Fatal(err)
		}

		expected := "Missing dependency: backend"
		_, merr := runner.get_dependency("nginx")
		assert.Equal(t, expected, merr.Error())
	})
}
