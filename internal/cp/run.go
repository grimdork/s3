package cp

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grimdork/climate/arg"
	"github.com/grimdork/climate/paths"
	"github.com/grimdork/s3/internal/cfg"
	"github.com/grimdork/s3/internal/download"
)

// Run the command.
func Run(opt *arg.Options) error {
	profile := "default"
	if cfg.Config.Default != "" {
		profile = cfg.Config.Default
	}

	// Flags
	o := arg.New("s3 cp")
	o.SetDefaultHelp(true)
	o.SetFlag(cfg.GroupFlags, "r", "recursive", "Copy recursively.")
	o.SetFlag(cfg.GroupFlags, "L", "list", "List of files to copy from the source. Ignores -r.")
	o.SetFlag(cfg.GroupFlags, "f", "force", "Force overwrite.")
	o.SetOption("Config", "p", "profile", "Profile to load.", profile, false, arg.VarString, nil)
	o.SetOption("Config", "u", "url", "URL for the API.", cfg.Config.Profiles[profile], false, arg.VarString, nil)

	// Source and destination
	o.SetPositional("SOURCE", "File, directory, bucket or object URI to copy from.", "", true, arg.VarString)
	o.SetPositional("DESTINATION", "File, directory, bucket or object URI to copy to.", "", true, arg.VarString)
	err := o.Parse(opt.Args)
	if err == arg.ErrNoArgs {
		println("No arguments provided.")
		opt.PrintHelp()
		os.Exit(1)
	}

	if err != nil {
		return err
	}

	srcurl := o.GetPosString("SOURCE")
	desturl := o.GetPosString("DESTINATION")
	// list := o.GetString("List")
	force := o.GetBool("f")

	if o.GetBool("r") {
		return arg.ErrRunCommand
	}

	src, size, err := download.OpenSource(srcurl)
	if err != nil {
		fmt.Printf("Error opening source: %s\n", err.Error())
		os.Exit(2)
	}

	st := download.IdentifyURI(srcurl)
	dt := download.IdentifyURI(desturl)

	var s3c *s3.Client
	if st == download.URIS3 || dt == download.URIS3 {
		s3c, err = cfg.S3Client(profile, o.GetString("url"))
		if err != nil {
			fmt.Printf("Error loading credentials: %s\n", err.Error())
			return arg.ErrRunCommand
		}

	}

	if dt == download.URILocal && paths.FileExists(desturl) && !force {
		fmt.Printf("Destination file exists. Use -f to overwrite.\n")
		src.Close()
		os.Exit(2)
	}

	if dt == download.URIS3 && !force {
		goa := &s3.GetObjectAttributesInput{
			Key: aws.String(desturl[4:]),
		}
		_, err = s3c.GetObjectAttributes(context.TODO(), goa)
		if err == nil {
			fmt.Printf("Destination file exists. Use -f to overwrite.\n")
			src.Close()
			os.Exit(2)
		}
	}

	dest, err := download.OpenDestination(desturl)
	if err != nil {
		fmt.Printf("Error opening destination: %s\n", err.Error())
		src.Close()
		os.Exit(2)
	}

	err = download.Copy(src, dest, size)
	println()
	if err != nil {
		fmt.Printf("Error copying: %s\n", err.Error())
		os.Exit(2)
	}

	return arg.ErrRunCommand
}
