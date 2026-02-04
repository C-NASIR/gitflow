package ui

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Table renders tabular output.
type Table struct {
	w *tabwriter.Writer
}

// NewTable constructs a table writer.
func NewTable(out io.Writer) *Table {
	return &Table{
		w: tabwriter.NewWriter(out, 0, 0, 2, ' ', 0),
	}
}

// Header writes a header row.
func (t *Table) Header(cols ...string) {
	for i, c := range cols {
		if i == len(cols)-1 {
			fmt.Fprintf(t.w, "%s\n", c)
		} else {
			fmt.Fprintf(t.w, "%s\t", c)
		}
	}
}

// Row writes a data row.
func (t *Table) Row(cols ...any) {
	for i, c := range cols {
		if i == len(cols)-1 {
			fmt.Fprintf(t.w, "%v\n", c)
		} else {
			fmt.Fprintf(t.w, "%v\t", c)
		}
	}
}

// KeyValue writes a two-column key/value row.
func (t *Table) KeyValue(key string, value any) {
	t.Row(key, value)
}

// Flush writes buffered table output.
func (t *Table) Flush() {
	_ = t.w.Flush()
}
