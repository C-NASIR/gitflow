package ui

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type Table struct {
	w *tabwriter.Writer
}

func NewTable(out io.Writer) *Table {
	return &Table{
		w: tabwriter.NewWriter(out, 0, 0, 2, ' ', 0),
	}
}

func (t *Table) Header(cols ...string) {
	for i, c := range cols {
		if i == len(cols)-1 {
			fmt.Fprintf(t.w, "%s\n", c)
		} else {
			fmt.Fprintf(t.w, "%s\t", c)
		}
	}
}

func (t *Table) Row(cols ...any) {
	for i, c := range cols {
		if i == len(cols)-1 {
			fmt.Fprintf(t.w, "%v\n", c)
		} else {
			fmt.Fprintf(t.w, "%v\t", c)
		}
	}
}

func (t *Table) Flush() {
	_ = t.w.Flush()
}
