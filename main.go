// Package main provides the gitflow CLI entrypoint.
package main

import "gitflow/cmd/root"

// main executes the root CLI command.
func main() {
	root.Execute()
}
