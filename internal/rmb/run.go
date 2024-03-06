package rmb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grimdork/climate/arg"
	"github.com/grimdork/s3/internal/cfg"
)

// Run the command.
func Run(opt *arg.Options) error {
	profile := "default"
	if cfg.Config.Default != "" {
		profile = cfg.Config.Default
	}

	// Flags
	o := arg.New("s3 rmb")
	o.SetDefaultHelp(true)
	o.SetOption("Config", "p", "profile", "Profile to load.", "default", false, arg.VarString, nil)
	o.SetOption("Config", "u", "url", "URL for the API.", cfg.Config.Profiles[profile], false, arg.VarString, nil)
	o.SetPositional("NAME", "Names of buckets to delete.", "", true, arg.VarStringSlice)
	err := o.Parse(opt.Args)
	if err == arg.ErrNoArgs {
		println("No arguments provided.")
		o.PrintHelp()
		return arg.ErrRunCommand
	}

	s3c, err := cfg.S3Client(profile, o.GetString("url"))
	if err != nil {
		fmt.Printf("Error loading credentials: %s\n", err.Error())
		return arg.ErrRunCommand
	}

	names := o.GetPosStringSlice("NAME")
	if len(names) == 0 {
		fmt.Println("No bucket names provided.")
		o.PrintHelp()
		return arg.ErrRunCommand
	}

	for _, name := range names {
		_, err = s3c.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
			Bucket: aws.String(name)})
		if err != nil {
			fmt.Printf("Error deleting bucket %s: %s\n", name, err.Error())
			return arg.ErrRunCommand
		}

		fmt.Printf("Deleted bucket %s\n", name)
	}

	return arg.ErrRunCommand
}
