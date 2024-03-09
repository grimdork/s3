package cp

import (
	"fmt"
	"os"

	"github.com/grimdork/climate/arg"
	"github.com/grimdork/s3/internal/cfg"
	"github.com/grimdork/s3/internal/download"
)

const groupFlags = "Flags"

// Run the command.
func Run(opt *arg.Options) error {
	profile := "default"
	if cfg.Config.Default != "" {
		profile = cfg.Config.Default
	}

	// Flags
	o := arg.New("s3 cp")
	o.SetDefaultHelp(true)
	o.SetFlag(groupFlags, "r", "recursive", "Copy recursively.")
	o.SetFlag(groupFlags, "L", "list", "List of files to copy from the source. Ignores -r.")
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
	// recursive := o.GetBool("r")
	// list := o.GetString("List")

	src, size, err := download.OpenSource(srcurl)
	if err != nil {
		fmt.Printf("Error opening source: %s\n", err.Error())
		os.Exit(2)
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

	// s3c, err := cfg.S3Client(profile, o.GetString("url"))
	// if err != nil {
	// 	fmt.Printf("Error loading credentials: %s\n", err.Error())
	// 	return arg.ErrRunCommand
	// }

	return arg.ErrRunCommand
}
