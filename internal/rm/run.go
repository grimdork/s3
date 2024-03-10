package rm

import (
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
	o := arg.New("s3 cp")
	o.SetDefaultHelp(true)
	o.SetFlag(cfg.GroupFlags, "r", "recursive", "Delete recursively.")
	o.SetFlag(cfg.GroupFlags, "L", "list", "List of files to delete. Ignores -r.")
	o.SetFlag(cfg.GroupFlags, "f", "force", "Don't ask for confirmation.")
	o.SetOption("Config", "p", "profile", "Profile to load.", profile, false, arg.VarString, nil)
	o.SetOption("Config", "u", "url", "URL for the API.", cfg.Config.Profiles[profile], false, arg.VarString, nil)

	return nil
}
