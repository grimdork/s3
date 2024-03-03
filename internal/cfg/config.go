package cfg

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/grimdork/climate/paths"
)

// Configuration for s3.
type Configuration struct {
	// Default profile to use when the -p option is not set.
	Default string `json:"default"`
	// Profiles to URLs mapping to make non-AWS storage easier.
	Profiles map[string]string `json:"profiles"`
}

// Config is the program configuration.
var Config Configuration

const (
	// ProgramName for use wherever needed.
	ProgramName = "s3"
	// ConfigPath to store settings in.
	ConfigPath = "net.grimdork.s3"
)

func init() {
	err := load()
	if err != nil {
		panic(err)
	}
}

// load the configuration.
func load() error {
	cfgpath, err := paths.New(ConfigPath)
	if err != nil {
		panic(err)
	}

	fn := filepath.Join(cfgpath.UserBase, "config.json")
	data, err := os.ReadFile(fn)
	if err != nil {
		Config = Configuration{Profiles: map[string]string{}}
		return nil
	}

	return json.Unmarshal(data, &Config)
}

// Save saves the configuration.
func Save() error {
	cfgpath, err := paths.New(ConfigPath)
	if err != nil {
		return err
	}

	if !paths.DirExists(cfgpath.UserBase) {
		err = os.MkdirAll(cfgpath.UserBase, 0700)
		if err != nil {
			return err
		}
	}

	fn := filepath.Join(cfgpath.UserBase, "config.json")
	data, err := json.MarshalIndent(Config, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(fn, data, 0600)
}
