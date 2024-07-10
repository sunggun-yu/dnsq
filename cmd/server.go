/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/dnsq/internal/server"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A sub command for server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Execute()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
