package config

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin"
)

const envNameTarget = "GIT_STATUSES_TARGET"

type Config struct {
	TargetPaths []string
}

func FromArgs(args []string, app *kingpin.Application) (*Config, error) {
	conf := Config{TargetPaths: []string{}}

	app.Arg("pathspec", "path for directory to find repositories").ExistingDirsVar(&conf.TargetPaths)
	if _, err := app.Parse(args); err != nil {
		return nil, err
	}

	if len(conf.TargetPaths) == 0 {
		conf.TargetPaths = filepath.SplitList(os.Getenv(envNameTarget))
	}

	if len(conf.TargetPaths) == 0 {
		conf.TargetPaths = []string{"."}
	}

	return &conf, nil
}
