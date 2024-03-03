package lp

import (
	"fmt"

	"github.com/grimdork/climate/arg"
	"github.com/grimdork/s3/internal/cfg"
)

// Run the command.
func Run(opt *arg.Options) error {
	fmt.Printf("%d existing profiles.\n", len(cfg.Config.Profiles))
	for k, v := range cfg.Config.Profiles {
		fmt.Printf("\t%s = %s\n", k, v)
	}

	if cfg.Config.Default != "" {
		fmt.Printf("Default profile: %s\n", cfg.Config.Default)
	}

	return arg.ErrRunCommand
}
