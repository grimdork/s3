package main

import (
	"fmt"

	"github.com/grimdork/climate/arg"
)

// VersionRun simply displays the version of the program.
func VersionRun(opt *arg.Options) error {
	fmt.Printf("%s version %s (%s)\n", appname, version, date)
	return arg.ErrRunCommand
}
