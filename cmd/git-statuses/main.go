package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	statuses "github.com/kyoh86/git-statuses"
	"github.com/saracen/walker"
	"github.com/spf13/cobra"
)

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

const (
	appName = "git-statuses"
)

var (
	formatter statuses.Formatter
	json      bool
)

var facadeCommand = &cobra.Command{
	Use:     appName,
	Short:   "git-statuses finds local git repositories and show statuses of them.",
	Version: fmt.Sprintf("%s-%s (%s)", version, commit, date),
	RunE: func(command *cobra.Command, targets []string) error {
		ctx := command.Context()
		if len(targets) == 0 {
			targets = []string{"."}
		}
		formatter = statuses.ShortFormat
		if json {
			formatter = statuses.JSONFormat
		}
		for _, target := range targets {
			return walker.WalkWithContext(ctx, target, processOne(ctx))
		}
		return nil
	},
}

func init() {
	facadeCommand.Flags().BoolVarP(&json, "json", "", false, "Format as JSON")
}

func processOne(ctx context.Context) func(path string, info os.FileInfo) error {
	return func(path string, info os.FileInfo) error {
		if !info.IsDir() || info.Name() != ".git" {
			return nil
		}

		repositoryPath := filepath.Dir(path)
		logger := log.FromContext(ctx).WithField("path", repositoryPath)
		ctx := log.NewContext(ctx, logger)

		state, err := statuses.GetStatus(ctx, repositoryPath)
		if err != nil {
			logger.Error(err.Error())
		}

		text, err := formatter(state)
		if err != nil {
			return err
		}
		if text != "" {
			fmt.Println(text)
		}
		return filepath.SkipDir // .git ディレクトリの下は見ない
	}
}

func main() {
	ctx := log.NewContext(context.Background(), &log.Logger{
		Handler: cli.New(os.Stderr),
		Level:   log.InfoLevel,
	})
	if err := facadeCommand.ExecuteContext(ctx); err != nil {
		log.FromContext(ctx).Error(err.Error())
		os.Exit(1)
	}
}
