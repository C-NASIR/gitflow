package ui

import (
	"fmt"
	"io"
)

type UI struct {
	Out     io.Writer
	Err     io.Writer
	Color   bool
	Emoji   bool
	Verbose bool
}

func New(out io.Writer, err io.Writer, color, emoji, verbose bool) *UI {
	return &UI{
		Out:     out,
		Err:     err,
		Color:   color,
		Emoji:   emoji,
		Verbose: verbose,
	}
}

func (u *UI) Header(text string) {
	fmt.Fprintln(u.Out)
	fmt.Fprintln(u.Out, text)
	fmt.Fprintln(u.Out)
}

func (u *UI) Line(format string, args ...any) {
	fmt.Fprintf(u.Out, format+"\n", args...)
}

func (u *UI) Success(format string, args ...any) {
	prefix := "OK "
	if u.Emoji {
		prefix = "✅ "
	}
	fmt.Fprintf(u.Out, prefix+format+"\n", args...)
}

func (u *UI) Warn(format string, args ...any) {
	prefix := "WARN "
	if u.Emoji {
		prefix = "⚠️ "
	}
	fmt.Fprintf(u.Out, prefix+format+"\n", args...)
}

func (u *UI) Error(format string, args ...any) {
	prefix := "ERR "
	if u.Emoji {
		prefix = "❌ "
	}
	fmt.Fprintf(u.Err, prefix+format+"\n", args...)
}
