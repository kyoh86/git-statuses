package statuses

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/apex/log"
)

func GetStatus(ctx context.Context, path string) (Status, error) {
	collector := &collector{ctx: ctx}
	cmd := exec.Command("git", "status", "--porcelain", "--branch", "--ahead-behind")
	cmd.Stdout = collector
	cmd.Stderr = os.Stderr
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return Status{}, err
	}
	if err := collector.Close(); err != nil {
		return Status{}, err
	}
	status := collector.Status
	remote, err := getRemote(path)
	if err != nil {
		return Status{}, err
	}
	url, err := getRemoteURL(path, remote)
	if err != nil {
		return Status{}, err
	}
	status.Path = path
	status.URL = url
	return status, nil
}

func getRemote(path string) (string, error) {
	cmd := exec.Command("git", "remote")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git remote: %w", err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func getRemoteURL(path, remote string) (string, error) {
	if remote == "" {
		return "", nil
	}
	cmd := exec.Command("git", "config", fmt.Sprintf("remote.%s.url", remote))
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git config: %w", err)
	}
	return strings.TrimSuffix(strings.TrimRight(string(out), "\n"), ".git"), nil
}

type Status struct {
	Path      string `json:"path"`
	Untracked bool   `json:"untracked"`
	Modified  bool   `json:"modified"`
	Branch    string `json:"branch"`
	Upstream  string `json:"upstream"`
	Ahead     int    `json:"ahead"`
	Behind    int    `json:"behind"`
	URL       string `json:"url"`
}

type collector struct {
	Status

	ctx     context.Context
	surplus string
	suspend bool
}

const (
	branchPrefix     = "## "
	branchInitPrefix = branchPrefix + "No commits yet on "
)

var (
	branchRegexp = regexp.MustCompile(`^## (\S+)\.\.\.(\S+/\S+)(?: \[(?:ahead (\d+))?(?:, )?(?:behind (\d+))?\])?$`)
	errStop      = errors.New("stop")
)

func parseInt32(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	i64, err := strconv.ParseInt(str, 10, 32)
	return int(i64), err
}

func (c *collector) parseStatusRune(r rune) {
	switch r {
	case '?':
		c.Untracked = true
	default:
		c.Modified = true
	}
}

func (c *collector) parseStatusLine(line string) error {
	if c.Untracked && c.Modified {
		c.suspend = true
	}
	if c.suspend {
		return nil
	}
	if len(line) < 2 {
		return nil
	}

	if !strings.HasPrefix(line, branchPrefix) {
		c.parseStatusRune([]rune(line)[0])
		c.parseStatusRune([]rune(line)[1])
		return nil
	}

	if strings.HasPrefix(line, branchInitPrefix) {
		c.Branch = "<none>"
		c.Upstream = "<none>"
		c.Untracked = true
		c.suspend = true
		return nil
	}
	for i, part := range branchRegexp.FindStringSubmatch(line) {
		switch i {
		case 4:
			behind, err := parseInt32(part)
			if err != nil {
				return fmt.Errorf("parse behind: %w", err)
			}
			c.Behind = behind
		case 3:
			ahead, err := parseInt32(part)
			if err != nil {
				return fmt.Errorf("parse ahead: %w", err)
			}
			c.Ahead = ahead
		case 2:
			c.Upstream = part
		case 1:
			c.Branch = part
		}
	}
	return nil
}

func (c *collector) Write(p []byte) (n int, _ error) {
	n = len(p)
	lines := strings.Split(c.surplus+string(p), "\n")
	last := len(lines) - 1
	c.surplus = lines[last]
	for i := 0; i < last; i++ {
		err := c.parseStatusLine(lines[i])
		switch {
		case err == nil:
			// noop
		case errors.Is(err, errStop):
			return
		default:
			log.FromContext(c.ctx).Error(err.Error())
		}
	}
	return
}

func (c *collector) Close() error {
	return c.parseStatusLine(c.surplus)
}
