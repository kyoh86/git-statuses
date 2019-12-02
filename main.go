package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/git-statuses/config"
	"github.com/kyoh86/git-statuses/git/local"
)

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	app := kingpin.New("git-statuses", "git-statuses finds local git repositories and show statuses of them.").Version(version).Author("kyoh86")
	conf, err := config.FromArgs(os.Args[1:], app)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	errorMap := map[string]error{}
	for _, targetPath := range conf.TargetPaths {
		targetPath := targetPath
		if err := local.WalkOnRepositories(targetPath, func(repositoryPath string) error {
			output := conf.WrapStatusOutput(targetPath, repositoryPath, os.Stdout)
			defer output.Close()

			errput := conf.WrapStatusOutput(targetPath, repositoryPath, os.Stderr)
			defer errput.Close()

			if err := conf.Status(repositoryPath, output, errput); err != nil {
				fmt.Fprintln(errput, err)
			}

			return nil
		}); err != nil {
			errorMap[targetPath] = err
		}
	}

	if len(errorMap) > 0 {
		fmt.Fprintln(os.Stderr, errorMap)
	}
}
