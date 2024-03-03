package sdp

import (
	"github.com/grimdork/climate/arg"
	"github.com/grimdork/s3/internal/cfg"
)

// Run the command.
func Run(opt *arg.Options) error {
	o := arg.New("s3 sdp")
	o.SetDefaultHelp(true)
	o.SetPositional("PROFILE", "Name of an S3 profile in the AWS config.", "", true, arg.VarString)
	err := o.Parse(opt.Args)
	if err != nil {
		return err
	}

	p := o.GetPosString("PROFILE")
	if p == "" {
		println("Error: ", arg.ErrMissingRequired.Error())
		o.PrintHelp()
		return arg.ErrRunCommand
	}

	cfg.Config.Default = p
	err = cfg.Save()
	if err != nil {
		return err
	}

	return arg.ErrRunCommand
}
