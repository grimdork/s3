package ls

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grimdork/climate/arg"
	"github.com/grimdork/climate/human"
	"github.com/grimdork/s3/internal/cfg"
)

// Run the command.
func Run(opt *arg.Options) error {
	profile := "default"
	if cfg.Config.Default != "" {
		profile = cfg.Config.Default
	}

	// Flags
	o := arg.New("s3 ls")
	o.SetDefaultHelp(true)
	o.SetOption("Config", "p", "profile", "Profile to load.", profile, false, arg.VarString, nil)
	o.SetOption("Config", "u", "url", "URL for the API.", cfg.Config.Profiles[profile], false, arg.VarString, nil)
	o.SetPositional("PATH", "Bucket/object path to list, or blank to list buckets.", "", false, arg.VarStringSlice)
	err := o.Parse(opt.Args)
	if err != nil {
		return err
	}

	s3c, err := cfg.S3Client(profile, o.GetString("url"))
	if err != nil {
		fmt.Printf("Error loading credentials: %s\n", err.Error())
		return arg.ErrRunCommand
	}

	// List S3 compatible buckets using selected profile.
	pathlist := o.GetPosStringSlice("PATH")
	if len(pathlist) == 0 {
		result, err := s3c.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		if err != nil {
			fmt.Printf("Couldn't list buckets for your account: %s\n", err.Error())
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

	for _, path := range pathlist {
		fmt.Printf("%s:\n", path)
		result, err := s3c.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: aws.String(path),
		})
		if err != nil {
			fmt.Printf("Couldn't list objects in bucket: %s\n", err.Error())
			return arg.ErrRunCommand
		}

		if len(result.Contents) == 0 {
			fmt.Println("\tNo objects found.")
		}

		for _, obj := range result.Contents {
			fmt.Printf("\t%s\t%s\n", *obj.Key, human.UInt(uint64(*obj.Size), false))
		}
	}

	return arg.ErrRunCommand
}
