package cfg

import (
	"errors"
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v2"
)

// Config contains configurations used on a per repo basis. Serializes to YAML.
type Config struct {
	Vcs       string    `yaml:"vcs"`       // What VCS are you using?
	Author    *Author   `yaml:"author"`    // Who's machine is this?
	Teammates []*Author `yaml:"teammates"` // Who's working with you?
	Path      string    // Where this config came from
}

// Author describes a project collaborator. Serialized to YAML.
type Author struct {
	Name  string `yaml:"name"`  // Author name. e.g. Lindsey Bluth
	Alias string `yaml:"alias"` // Nickname. e.g. lb
	Email string `yaml:"email"` // Email address. e.g. lindsb@example.com
}

// ByName implements sort.Interface for []*Author based on the author name.
type ByName []*Author

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	return a[i].Name+a[i].Alias < a[j].Name+a[j].Alias
}

// New creates a new Config which will be located at the specified path when
// it's saved.
func New(path string) *Config {
	return &Config{Path: path}
}

// NewFromFile creates a new Config from the file located at the specified path.
func NewFromFile(path string) (*Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := Config{Path: path}
	if err := yaml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Reload reloads the config information from the path on disk.
func (c *Config) Reload() error {
	updated, err := NewFromFile(c.Path)
	if err != nil {
		return err
	}
	c.Vcs = updated.Vcs
	c.Author = updated.Author
	c.Teammates = updated.Teammates
	return nil
}

// Save saves the config to disk.
func (c *Config) Save() error {
	buf, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.Path, buf, 0644)
}

// Validate checks that an in-memory configuration is ok.
func (c *Config) Validate() (bool, error) {
	if c.Vcs == "" {
		return false, errors.New("vcs can't be empty")
	}
	if c.Author == nil {
		return false, errors.New("author can't be nil")
	}
	if c.Author.Email == "" {
		return false, errors.New("author.email is required")
	}
	return true, nil
}

func (c *Config) equals(other *Config) bool {
	if c == other {
		return true
	}
	if c.Path != other.Path {
		return false
	}
	if c.Vcs != other.Vcs {
		return false
	}
	if *c.Author != *other.Author {
		return false
	}
	if len(c.Teammates) != len(other.Teammates) {
		return false
	}
	sort.Sort(ByName(c.Teammates))
	sort.Sort(ByName(other.Teammates))
	for i := 0; i < len(c.Teammates); i++ {
		if *c.Teammates[i] != *other.Teammates[i] {
			print("here")
			return false
		}
	}
	return true
}
