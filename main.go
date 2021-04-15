package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

func main() {
	ctx := log.NewContext(context.Background(), &log.Logger{
		Handler: cli.New(os.Stderr),
		Level:   log.InfoLevel,
	})
	facadeCommand.Flags().BoolVarP(&flags.Detail, "detail", "d", false, "show detail results")
	facadeCommand.Flags().BoolVarP(&flags.Relative, "relative", "r", false, "show relative results")
	if err := facadeCommand.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

var flags struct {
	Detail   bool
	Relative bool
}

var facadeCommand = &cobra.Command{
	Use:     "git-statuses",
	Short:   "git-statuses finds local git repositories and show statuses of them",
	Version: fmt.Sprintf("%s-%s (%s)", version, commit, date),
	RunE: func(cmd *cobra.Command, paths []string) error {
		if len(paths) == 0 {
			paths = []string{"."}
		}
		errorMap := map[string]error{}
		ctx := context.Background()
		for _, targetPath := range paths {
			if err := filepath.Walk(targetPath, func(children string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() || info.Name() != ".git" {
					return nil
				}

				path := filepath.Dir(filepath.Join(targetPath, children))
				ctx = log.NewContext(ctx, log.FromContext(ctx).WithField("path", path))
				if err := statDir(ctx, path); err != nil {
					log.FromContext(ctx).WithField("error", err).Info("failed to get stat")
				}
				return filepath.SkipDir
			}); err != nil {
				errorMap[targetPath] = err
			}
		}
		return nil
	},
}

func statDir(ctx context.Context, path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("open a repository: %w", err)
	}
	var (
		ahead    int
		behind   int
		modified bool
		untrack  bool
	)
	tree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("get worktree: %w", err)
	}
	stat, err := tree.Status()
	if err != nil {
		// return fmt.Errorf("get status: %w", err)
	}
	for _, file := range stat {
		if file.Staging == git.Unmodified && file.Worktree == git.Unmodified {
		} else if file.Staging == git.Untracked {
			untrack = true
		} else {
			modified = true
		}
	}
	localHead, err := repo.Head()
	if err != nil {
		return fmt.Errorf("get a HEAD refer: %w", err)
	}
	br, err := repo.Branch(localHead.Name().Short())
	if err != nil {
		return fmt.Errorf("get a HEAD branch: %w", err)
	}
	if br.Remote == "" {
		//TODO: untracking branch
	} else {
		// remote
		remoteHead, err := repo.Reference(plumbing.NewRemoteReferenceName(br.Remote, br.Name), true)
		if err != nil {
			return fmt.Errorf("get a remote HEAD: %w", err)
		}
		ahead, err = countCommit(ctx, repo, localHead, remoteHead)
		if err != nil {
			return fmt.Errorf("count ahead: %w", err)
		}
		behind, err = countCommit(ctx, repo, remoteHead, localHead)
		if err != nil {
			return fmt.Errorf("count behind: %w", err)
		}
	}
	fmt.Println(ahead, behind, modified, untrack)
	return nil
}

func countCommit(ctx context.Context, repo *git.Repository, until *plumbing.Reference, since *plumbing.Reference) (int, error) {
	commit, err := repo.CommitObject(since.Hash())
	if err != nil {
		return 0, err
	}
	logs, err := repo.Log(&git.LogOptions{
		From:  until.Hash(),
		Since: &commit.Author.When,
	})
	if err != nil {
		return 0, err
	}
	inc := 1
	cnt := 0
	if err := logs.ForEach(func(c *object.Commit) error {
		switch c.Hash {
		case until.Hash():
			inc = 1
		case since.Hash():
			cnt += inc
			inc = 0
		default:
			cnt += inc
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return cnt, nil
}
