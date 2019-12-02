package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/git-statuses/git/local"
	"github.com/kyoh86/git-statuses/util"
)

const envNameTarget = "GIT_STATUSES_TARGET"

type Config struct {
	TargetPaths []string
	Detail      bool
	Relative    bool
}

func (c *Config) WrapStatusOutput(rootPath string, path string, base io.Writer) io.WriteCloser {
	var label string
	if c.Relative {
		label, _ = filepath.Rel(rootPath, path)
	} else {
		label = path
	}
	if c.Detail {
		return util.NewLabeledWriter(label, base)
	}
	return util.NewSimpleLabeledWriter(label, base)
}

func (c *Config) Status(path string, out io.Writer, err io.Writer) error {
	if c.Detail {
		return local.Status(path, out, err)
	}
	return local.ShortStatus(path, out, err)
}

func FromArgs(args []string, app *kingpin.Application) (*Config, error) {
	conf := Config{TargetPaths: []string{}}

	app.Arg("pathspec", "path for directory to find repositories").ExistingDirsVar(&conf.TargetPaths)
	app.Flag("detail", "show detail results").Short('d').BoolVar(&conf.Detail)
	app.Flag("relative", "show relative results").Short('r').BoolVar(&conf.Relative)
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
