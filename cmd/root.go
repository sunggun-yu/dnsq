package cmd

import (
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/dnsq/internal/dnslookup"
	"github.com/sunggun-yu/dnsq/internal/models"
)

// root command
var rootCmd = rootCommand()

// rootCmd represents the base command when called without any subcommands
func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "dnsq [domains...]",
		Short:         "Look up DNS records for one or more domains",
		SilenceErrors: true,
		Args:          cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			records := dnslookup.GetDNSRecords(args)
			printRecords(cmd.OutOrStdout(), records)
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

// Set the version of cmd
func SetVersion(version string) {
	rootCmd.Version = version
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
func printRecords(w io.Writer, records map[string][]models.DNSRecord) {

	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.Style().Options.DrawBorder = true
	tw.Style().Options.SeparateHeader = true
	tw.Style().Options.SeparateRows = false
	tw.Style().Options.SeparateColumns = true

	tw.AppendHeader(table.Row{"Domain", "Host", "Type", "Data"})
	for domain, data := range records {
		for _, d := range data {
			tw.AppendRow(table.Row{domain, d.Host, d.Type, d.Data})
		}
	}
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	tw.Style().Options.SeparateRows = true
	w.Write([]byte(tw.Render()))
	w.Write([]byte("\n"))
}
