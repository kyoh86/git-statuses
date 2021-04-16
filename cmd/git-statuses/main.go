package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	statuses "github.com/kyoh86/git-statuses"
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
	RunE: func(_ *cobra.Command, targets []string) error {
		if len(targets) == 0 {
			targets = []string{"."}
		}
		formatter = statuses.ShortFormat
		if json {
			formatter = statuses.JSONFormat
		}
		for _, target := range targets {
			return filepath.Walk(target, processOne)
		}
		return nil
	},
}

func init() {
	facadeCommand.Flags().BoolVarP(&json, "json", "", false, "Format as JSON")
}

func processOne(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Name() != ".git" {
		return nil
	}

	repositoryPath := filepath.Dir(path)
	state, err := statuses.GetStatus(repositoryPath)
	if err != nil {
		log.Error(err.Error())
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

func main() {
	ctx := log.NewContext(context.Background(), &log.Logger{
		Handler: cli.New(os.Stderr),
		Level:   log.InfoLevel,
	})
	if err := facadeCommand.ExecuteContext(ctx); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
