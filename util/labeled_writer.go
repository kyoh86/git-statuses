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

func (w *labeledWriter) Write(p []byte) (int, error) {
	str := string(p)
	w.once.Do(func() {
		w.writeCore(w.label)
		w.writeCore(newline)
	})

	w.writeCore(w.prefix)
	w.writeCore(strings.Replace(strings.TrimSuffix(str, newline), newline, w.replace, -1))

	if strings.HasSuffix(str, newline) {
		w.writeCore(newline)
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
