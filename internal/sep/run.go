package sep

import (
	"github.com/grimdork/climate/arg"
	"github.com/grimdork/s3/internal/cfg"
)

// Run the command.
func Run(opt *arg.Options) error {
	o := arg.New("s3 sep")
	o.SetDefaultHelp(true)
	o.SetPositional("PROFILE", "Name of an S3 profile in the AWS config.", "", true, arg.VarString)
	o.SetPositional("ENTRYPOINT", "Base path for the S3 API entrypoint.", "", true, arg.VarString)
	err := o.Parse(opt.Args)
	if err != nil {
		println("FAIL: ", err.Error())
		return err
	}

	profile := o.GetPosString("PROFILE")
	ep := o.GetPosString("ENTRYPOINT")
	if ep == "" {
		println("Error: ", arg.ErrMissingRequired.Error())
		o.PrintHelp()
		return arg.ErrRunCommand
	}

	cfg.Config.Profiles[profile] = ep
	err = cfg.Save()
	if err != nil {
		return err
	}

	return arg.ErrRunCommand
}
