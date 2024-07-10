package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunggun-yu/dnsq/internal/dnslookup"
	"github.com/sunggun-yu/dnsq/internal/models"
)

// root command
var rootCmd = rootCommand()

// rootCmd represents the base command when called without any subcommands
func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dlk [domains...]",
		Short: "Look up DNS records for one or more domains",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, domain := range args {
				records := dnslookup.GetDNSRecords(domain)
				printRecords(cmd.OutOrStdout(), records)
			}
		},
	}
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dlk.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// printRecords prints DNS records to the console
func printRecords(w io.Writer, records []models.DNSRecord) {
	for _, record := range records {
		w.Write([]byte(fmt.Sprintf("%s %s %s\n", record.Host, record.Type, record.Data)))
	}
}
