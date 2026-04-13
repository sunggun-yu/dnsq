package cmd

import (
	"embed"

	"github.com/spf13/cobra"
	"github.com/sunggun-yu/dnsq/internal/server"
)

// StaticFS is an embed.FS that contains static files from static directory
var StaticFS embed.FS

// flags struct for server command
type serverFlags struct {
	port int
}

// serverCmd represents the server command
var serverCmd = serverCommand()

// serverCommand returns a cobra.Command for server
func serverCommand() *cobra.Command {

	// add flags
	var flags serverFlags

	// create a new cobra.Command
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the DNS lookup web server",
		Run: func(cmd *cobra.Command, args []string) {
			srv := server.NewServer(flags.port, StaticFS, rootCmd.Version)
			srv.Run()
		},
	}

	// add flags
	cmd.Flags().IntVarP(&flags.port, "port", "p", 8080, "port number for the server")

	return cmd
}

// init registers the server command
func init() {
	rootCmd.AddCommand(serverCmd)
}
