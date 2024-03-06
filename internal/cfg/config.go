package cfg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

// S3Client builder.
func S3Client(profile, url string) (*s3.Client, error) {
	s3cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	if url != "" {
		fmt.Printf("Using profile %s with API entrypoint %s\n", profile, url)
		s3cfg.BaseEndpoint = aws.String(url)
	} else {
		println("Defaults")
		s3cfg.BaseEndpoint = aws.String("https://s3.eu-central-1.amazonaws.com")
		s3cfg.Region = "eu-central-1"
		cred, err := s3cfg.Credentials.Retrieve(context.Background())
		if err != nil {
			return nil, err
		}
		fmt.Printf("Access Key: %v\n", cred.AccessKeyID)
	}
	return s3.New(s3.Options{}, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = s3cfg.BaseEndpoint
		o.Region = s3cfg.Region
		o.Credentials = s3cfg.Credentials
	}), nil
}
