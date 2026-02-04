package ui

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

type UI struct {
	out      io.Writer
	useColor bool
	useEmoji bool
	verbose  bool

	okPrefix   string
	warnPrefix string
	errPrefix  string
	infoPrefix string

	okStyle   *color.Color
	warnStyle *color.Color
	errStyle  *color.Color
	infoStyle *color.Color
	headStyle *color.Color
}

type Options struct {
	Out     io.Writer
	Color   bool
	Emoji   bool
	Verbose bool
}

func New(opts Options) *UI {
	u := &UI{
		out:      opts.Out,
		useColor: opts.Color,
		useEmoji: opts.Emoji,
		verbose:  opts.Verbose,
	}

	if u.useEmoji {
		u.okPrefix = "✓ "
		u.warnPrefix = "⚠ "
		u.errPrefix = "✗ "
		u.infoPrefix = "• "
	} else {
		u.okPrefix = "OK "
		u.warnPrefix = "WARN "
		u.errPrefix = "ERR "
		u.infoPrefix = "INFO "
	}

	if u.useColor {
		u.okStyle = color.New(color.FgGreen, color.Bold)
		u.warnStyle = color.New(color.FgYellow, color.Bold)
		u.errStyle = color.New(color.FgRed, color.Bold)
		u.infoStyle = color.New(color.FgCyan)
		u.headStyle = color.New(color.Bold, color.Underline)
	} else {
		u.okStyle = color.New()
		u.warnStyle = color.New()
		u.errStyle = color.New()
		u.infoStyle = color.New()
		u.headStyle = color.New()
	}

	return u

}

func (u *UI) Header(text string) {
	fmt.Fprintln(u.out)
	u.headStyle.Fprintln(u.out, text)
	fmt.Fprintln(u.out)
}

func (u *UI) Line(format string, args ...any) {
	fmt.Fprintf(u.out, format+"\n", args...)
}

func (u *UI) Success(format string, args ...any) {
	u.okStyle.Fprintf(u.out, u.okPrefix+format+"\n", args...)
}

func (u *UI) Warn(format string, args ...any) {
	u.warnStyle.Fprintf(u.out, u.warnPrefix+format+"\n", args...)
}

func (u *UI) Error(format string, args ...any) {
	u.errStyle.Fprintf(u.out, u.errPrefix+format+"\n", args...)
}

func (u *UI) Info(format string, args ...any) {
	u.infoStyle.Fprintf(u.out, u.infoPrefix+format+"\n", args...)
}

func (u *UI) Verbose(format string, args ...any) {
	if !u.verbose {
		return
	}
	u.Info(format, args...)
}
