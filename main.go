package main

import (
	"fmt"
	"os"

	"github.com/grimdork/climate/arg"
	"github.com/grimdork/s3/internal/cp"
	"github.com/grimdork/s3/internal/lp"
	"github.com/grimdork/s3/internal/ls"
	"github.com/grimdork/s3/internal/mb"
	"github.com/grimdork/s3/internal/rmb"
	"github.com/grimdork/s3/internal/sdp"
	"github.com/grimdork/s3/sep"
)

const (
	groupCfg = "Configuration"
	groupS3  = "S3"
)

func main() {
	opt := arg.New("s3")
	opt.SetDefaultHelp(true)
	opt.SetOption(arg.GroupDefault, "v", "version", "Display the version and exit.", false, false, arg.VarBool, nil)
	opt.SetCommand("lp", "List S3 profile entrypoints.", groupCfg, lp.Run, []string{"listprofiles"})
	opt.SetCommand("sep", "Set the entrypoint for an S3 profile.", groupCfg, sep.Run, []string{"setentrypoint"})
	opt.SetCommand("sdp", "Set the default S3 profile to use.", groupCfg, sdp.Run, []string{"setdefaultprofile"})
	opt.SetCommand("mb", "Make buckets.", groupS3, mb.Run, nil)
	opt.SetCommand("rmb", "Remove buckets.", groupS3, rmb.Run, nil)
	opt.SetCommand("ls", "List buckets or contents of buckets.", groupS3, ls.Run, nil)
	opt.SetCommand("cp", "Copy files, folders, buckets or objects.", groupS3, cp.Run, nil)
	err := opt.Parse(os.Args)
	if err == arg.ErrRunCommand {
		return
	}

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		opt.PrintHelp()
		return
	}

	opt.PrintHelp()
}
