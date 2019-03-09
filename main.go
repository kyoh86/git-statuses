package main

import (
	"fmt"
	"os"

	"github.com/kyoh86/git-statuses/config"
	"github.com/kyoh86/git-statuses/git/local"
)

func main() {
	conf, err := config.FromArgs(os.Args[1:])
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
