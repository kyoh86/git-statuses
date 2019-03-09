package local

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

func Status(path string, out, err io.Writer) error {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Stdout = out
	cmd.Stderr = err
	cmd.Dir = path
	return cmd.Run()
}

type RepositoryStatus int

const (
	StatusCodeClear     = RepositoryStatus(0)
	StatusCodeModified  = RepositoryStatus(1)
	StatusCodeUntracked = RepositoryStatus(2)
)

func runeToCode(r rune) RepositoryStatus {
	switch r {
	case ' ':
		return StatusCodeClear
	case '?':
		return StatusCodeUntracked
	default:
		return StatusCodeModified
	}
}

func (c RepositoryStatus) String() string {
	switch c {
	default: //case StatusCodeClear:
		return "  "
	case StatusCodeUntracked:
		return " +"
	case StatusCodeModified:
		return "M "
	case StatusCodeUntracked | StatusCodeModified:
		return "M+"
	}
}

type RepositoryStatusCollector struct {
	Code RepositoryStatus

	surplus string
	once    sync.Once
}

func (w *RepositoryStatusCollector) init() {
	w.once.Do(func() {
		w.reset()
	})
}

func (w *RepositoryStatusCollector) reset() {
	w.Code = StatusCodeClear
	w.surplus = ""
}

func (w *RepositoryStatusCollector) parseLine(line string) {
	if len(line) > 2 {
		w.Code |= runeToCode([]rune(line)[0])
		w.Code |= runeToCode([]rune(line)[1])
	}
}

func (w *RepositoryStatusCollector) Write(p []byte) (int, error) {
	w.init()

	lines := strings.Split(w.surplus+string(p), "\n")
	last := len(lines) - 1
	for i := 0; i < last; i++ {
		w.parseLine(lines[i])
	}
	w.surplus = lines[last]
	return len(p), nil
}

func (w *RepositoryStatusCollector) Close() RepositoryStatus {
	w.init()

	w.parseLine(w.surplus)
	defer w.reset()
	return w.Code
}

func ShortStatus(path string, out io.Writer, err io.Writer) error {
	collector := &RepositoryStatusCollector{}
	if err := Status(path, collector, err); err != nil {
		return err
	}
	code := collector.Close()

	if code != StatusCodeClear {
		fmt.Fprintln(out, code.String())
	}
	return nil
}
