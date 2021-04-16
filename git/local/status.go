package local

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/apex/log"
)

func Status(path string, out io.Writer, err io.Writer) (RepositoryState, error) {
	collector := &RepositoryStatusCollector{}
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Stdout = collector
	cmd.Stderr = err
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return RepositoryState{}, err
	}
	collector.Close()

	return collector.RepositoryState, nil
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

type RepositoryState struct {
	Code     RepositoryStatus
	Branch   string
	Upstream string
	Ahead    int
	Behind   int
}

type RepositoryStatusCollector struct {
	RepositoryState

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

const (
	branchPrefix     = "## "
	branchInitPrefix = branchPrefix + "No commits yet on "
)

var (
	branchRegexp = regexp.MustCompile(`^## (\S+)\.\.\.(\S+/\S+)(?: \[(?:ahead (\d+))?(?:, )?(?:behind (\d+))?\])?$`)
)

func parseInt32(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	i64, err := strconv.ParseInt(str, 10, 32)
	return int(i64), err
}

func (w *RepositoryStatusCollector) parseLine(line string) {
	if err := w.parseLineCore(line); err != nil {
		log.Error(err.Error())
	}
}

func (w *RepositoryStatusCollector) parseLineCore(line string) error {
	if len(line) < 2 {
		return nil
	}

	if strings.HasPrefix(line, branchPrefix) {
		if strings.HasPrefix(line, branchInitPrefix) {
			return errors.New("initial branch")
		}
		matches := branchRegexp.FindStringSubmatch(line)
		switch len(matches) {
		default:
			//noop
		case 5:
			behind, err := parseInt32(matches[4])
			if err != nil {
				return fmt.Errorf("parse behind: %w", err)
			}
			w.Behind = behind
			fallthrough
		case 4:
			ahead, err := parseInt32(matches[3])
			if err != nil {
				return fmt.Errorf("parse ahead: %w", err)
			}
			w.Ahead = ahead
			fallthrough
		case 3:
			w.Upstream = matches[2]
			fallthrough
		case 2:
			w.Branch = matches[1]
		case 0, 1:
			//noop
		}

	} else {
		w.Code |= runeToCode([]rune(line)[0])
		w.Code |= runeToCode([]rune(line)[1])
	}
	return nil
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

func (w *RepositoryStatusCollector) Close() {
	w.init()

	w.parseLine(w.surplus)
	w.reset()
}
