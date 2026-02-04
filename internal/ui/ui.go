// Package ui provides CLI presentation helpers.
package ui

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// UI renders styled output to a writer.
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

// Options configures a UI instance.
type Options struct {
	Out     io.Writer
	Color   bool
	Emoji   bool
	Verbose bool
}

// New constructs a new UI using the supplied options.
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

// Header prints a formatted section header.
func (u *UI) Header(text string) {
	fmt.Fprintln(u.out)
	u.headStyle.Fprintln(u.out, text)
	fmt.Fprintln(u.out)
}

// Line writes a formatted line to output.
func (u *UI) Line(format string, args ...any) {
	fmt.Fprintf(u.out, format+"\n", args...)
}

// Success writes a success message.
func (u *UI) Success(format string, args ...any) {
	u.okStyle.Fprintf(u.out, u.okPrefix+format+"\n", args...)
}

// Warn writes a warning message.
func (u *UI) Warn(format string, args ...any) {
	u.warnStyle.Fprintf(u.out, u.warnPrefix+format+"\n", args...)
}

// Error writes an error message.
func (u *UI) Error(format string, args ...any) {
	u.errStyle.Fprintf(u.out, u.errPrefix+format+"\n", args...)
}

// Info writes an informational message.
func (u *UI) Info(format string, args ...any) {
	u.infoStyle.Fprintf(u.out, u.infoPrefix+format+"\n", args...)
}

// Verbose writes a message only when verbose output is enabled.
func (u *UI) Verbose(format string, args ...any) {
	if !u.verbose {
		return
	}
	u.Info(format, args...)
}

// StatusLabel formats a status label string.
func (u *UI) StatusLabel(level string) string {
	switch level {
	case "OK":
		return u.okStyle.Sprint(level)
	case "WARN":
		return u.warnStyle.Sprint(level)
	case "ERROR":
		return u.errStyle.Sprint(level)
	default:
		return level
	}
}

// ColorEnabled reports whether color output is enabled.
func (u *UI) ColorEnabled() bool {
	return u.useColor
}

// EmojiEnabled reports whether emoji output is enabled.
func (u *UI) EmojiEnabled() bool {
	return u.useEmoji
}

// VerboseEnabled reports whether verbose output is enabled.
func (u *UI) VerboseEnabled() bool {
	return u.verbose
}
