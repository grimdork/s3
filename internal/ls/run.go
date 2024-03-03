package ls

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grimdork/climate/arg"
	"github.com/grimdork/s3/internal/cfg"
)

// Run the command.
func Run(opt *arg.Options) error {
	// Flags
	o := arg.New("s3 ls")
	o.SetDefaultHelp(true)
	def := "default"
	if cfg.Config.Default != "" {
		def = cfg.Config.Default
	}
	o.SetOption("Config", "p", "profile", "Profile to load.", def, false, arg.VarString, nil)
	o.SetOption("Config", "u", "url", "URL for the API.", cfg.Config.Profiles[def], false, arg.VarString, nil)
	o.SetPositional("URI", "Location URIs to list, or blank to list buckets.", "", false, arg.VarStringSlice)
	err := o.Parse(opt.Args)
	if err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(o.GetString("profile")),
	)
	if err != nil {
		return err
	}

	base := o.GetString("url")
	if base != "" {
		cfg.BaseEndpoint = aws.String(base)
	} else {
		println("Defaults")
		cfg.BaseEndpoint = aws.String("https://s3.eu-central-1.amazonaws.com")
		cfg.Region = "eu-central-1"
		cred, err := cfg.Credentials.Retrieve(context.Background())
		if err != nil {
			return err
		}
		fmt.Printf("Access Key: %v\n", cred.AccessKeyID)
	}
	s3c := s3.NewFromConfig(cfg)

	// List S3 compatible buckets using selected profile.
	uris := o.GetPosStringSlice("URI")
	if len(uris) == 0 {
		result, err := s3c.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		if err != nil {
			fmt.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
			return arg.ErrRunCommand
		}
		if len(result.Buckets) == 0 {
			fmt.Println("You don't have any buckets.")
		} else {
			for _, bucket := range result.Buckets {
				fmt.Printf("\t%v\n", *bucket.Name)
			}
		}
		return arg.ErrRunCommand
	}

	return arg.ErrRunCommand
}
