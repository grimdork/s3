package cp

import (
	"os"

	"github.com/grimdork/climate/arg"
)

const groupFlags = "Flags"

// Run the command.
func Run(opt *arg.Options) error {
	// Flags
	o := arg.New("s3 cp")
	o.SetDefaultHelp(true)
	o.SetOption(groupFlags, "r", "recursive", "Copy recursively.", "", false, arg.VarString, nil)
	o.SetOption(groupFlags, "L", "list", "List of files to copy from.", "", false, arg.VarString, nil)
	o.SetOption("Config", "p", "profile", "Profile to load.", "default", false, arg.VarString, nil)

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

	return arg.ErrRunCommand
}
