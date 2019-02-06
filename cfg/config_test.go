package cfg

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	config      *Config
	testConfigs = map[string]string{
		"emptyConfig":   ``,
		"garbageFields": `abcd: 1234`,
		"sparseAuthor": `author:
  alias: mb`,
		"fullConfig": `vcs: git
author:
  name: Michael Bluth
  alias: mb
  email: mb@example.com
teammates:
  - name: Lindsey Bluth
    alias: lb
  - name: George Bluth
    alias: gb`,
	}
)

func TestNew(t *testing.T) {
	config = New("/tmp/cfg.yml")
	if config.Path != "/tmp/cfg.yml" {
		t.Fatalf("expected Path to set appropriately, was %v", config.Path)
	}
	if config.Vcs != "" {
		t.Fatalf("expected Vcs to be unset, got %v", config.Vcs)
	}
	if config.Author != nil {
		t.Fatalf("expected Author to be unset, got %v", config.Author)
	}
	if config.Teammates != nil {
		t.Fatalf("expected Teammates to be unset, got %v", config.Teammates)
	}
}

func TestNewFromFile(t *testing.T) {
	// test data

	writeFileContents := func(yml string) string {
		// set up
		f, err := ioutil.TempFile("", "config-*.yml")
		if err != nil {
			t.Fatalf("couldn't make tempfile during test set up: %v", err)
		}
		if _, err := f.Write([]byte(yml)); err != nil {
			t.Fatalf("couldn't write to tempfile during test set up: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("error closing file: %v", err)
		}
		return f.Name()
	}

	for id := range testConfigs {
		fname := writeFileContents(testConfigs[id])
		defer os.Remove(fname) // clean up
		config, err := NewFromFile(fname)
		if err != nil {
			t.Fatalf("error in NewFromFile: %v", err)
		}
		switch id {
		case "garbageFields":
			// if the config is just garbage, it should be same as empty
			fallthrough
		case "emptyConfig":
			if config.Vcs != "" || config.Author != nil || config.Teammates != nil {
				t.Fatalf("expected empty config to yeild no information")
			}
		case "sparseAuthor":
			if config.Author == nil {
				t.Fatalf("expected author to be created if some fields present")
			}
		case "fullConfig":
			if config.Vcs != "git" {
				t.Fatalf("expected vcs: git, got %v", config.Vcs)
			}
			compare := Author{
				Name:  "Michael Bluth",
				Alias: "mb",
				Email: "mb@example.com",
			}
			if *config.Author != compare {
				t.Fatalf("got unexpected author: %v", config.Author)
			}
		}
	}
}

func TestSave(t *testing.T) {
	f, _ := ioutil.TempFile("", "config-*.yml")
	defer os.Remove(f.Name()) // clean up
	config = &Config{
		Vcs:    "git",
		Author: &Author{},
		Teammates: []*Author{
			&Author{},
		},
		Path: f.Name(),
	}
	err := config.Save()
	if err != nil {
		t.Fatalf("error saving config: %v", err)
	}
	written, _ := NewFromFile(f.Name())
	if !config.equals(written) {
		t.Fatal("saved config was not equal to in memory config")
	}
}

func TestReload(t *testing.T) {
}

func TestValidate(t *testing.T) {

}
