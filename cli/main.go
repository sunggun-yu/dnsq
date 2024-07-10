/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/sunggun-yu/dnsq/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Version returns version and build information. it will be injected from ldflags(goreleaser)
func Version() string {
	return fmt.Sprintf("%s, commit %s, built at %s", version, commit, date)
}

func main() {
	cmd.SetVersion(Version())
	cmd.Execute()
}
