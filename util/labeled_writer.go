package util

import (
	"io"
	"strings"
	"sync"
)

type labeledWriter struct {
	label  string
	base   io.Writer
	indent string

	prefix  string
	replace string
	once    sync.Once
}

const (
	newline       = "\n"
	DefaultIndent = "  "
)

func NewLabeledWriter(label string, base io.Writer) io.WriteCloser {
	writer := &labeledWriter{
		base: base,
	}
	writer.SetIndent(DefaultIndent)
	writer.SetLabel(label)
	return writer
}

func (w *labeledWriter) SetIndent(indent string) {
	w.indent = indent
	w.replace = newline + indent
}

func (w *labeledWriter) SetLabel(label string) {
	w.label = label
	w.prefix = w.indent
	w.once = sync.Once{}
}

func (w *labeledWriter) Write(p []byte) (n int, retErr error) {
	str := string(p)
	w.once.Do(func() {
		if _, err := w.writeCore(w.label); err != nil && retErr == nil {
			retErr = err
			return
		}
		if _, err := w.writeCore(newline); err != nil && retErr == nil {
			retErr = err
		}
	})
	if retErr != nil {
		return 0, retErr
	}

	if _, err := w.writeCore(w.prefix); err != nil {
		return 0, err
	}
	if _, err := w.writeCore(strings.Replace(strings.TrimSuffix(str, newline), newline, w.replace, -1)); err != nil {
		return 0, err
	}

	if strings.HasSuffix(str, newline) {
		if _, err := w.writeCore(newline); err != nil {
			return 0, err
		}
		w.prefix = w.indent
	} else {
		w.prefix = ""
	}

	return len(p), nil
}

func (w *labeledWriter) Close() error {
	//FIXME: if w.base is io.Closer
	return nil
}

func (w *labeledWriter) writeCore(str string) (int, error) {
	if len(str) == 0 {
		return 0, nil
	}
	return w.base.Write([]byte(str))
}
