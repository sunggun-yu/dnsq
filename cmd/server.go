/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/dnsq/internal/server"
)

// flags struct for server command
type serverFlags struct {
	port int
}

// serverCmd represents the server command
var serverCmd = serverCommand()

func serverCommand() *cobra.Command {

	// add flags
	var flags serverFlags

	cmd := &cobra.Command{
		Use:   "server",
		Short: "A sub command for server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Run(flags.port)
		},
	}

	cmd.Flags().IntVarP(&flags.port, "port", "p", 8080, "port number for the server")

	return cmd
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
