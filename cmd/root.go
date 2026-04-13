package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/dnsq/internal/dnslookup"
	"github.com/sunggun-yu/dnsq/internal/models"
)

// root command
var rootCmd = rootCommand()

// rootCmd represents the base command when called without any subcommands
func rootCommand() *cobra.Command {

	var nameservers []string
	var includeDefault bool
	var ipv6 bool

	cmd := &cobra.Command{
		Use:           "dnsq [domains...]",
		Short:         "Look up DNS records for one or more domains",
		SilenceErrors: true,
		Args:          cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Build nameserver list
			finalNameservers := buildCLINameserverList(nameservers, includeDefault)
			results := dnslookup.GetDNSRecords(args, finalNameservers, ipv6)
			printRecords(cmd.OutOrStdout(), results)
		},
	}

	cmd.Flags().StringSliceVarP(&nameservers, "nameserver", "n", nil, "nameserver to query (can be specified multiple times)")
	cmd.Flags().BoolVarP(&includeDefault, "include-default", "d", false, "include default nameserver alongside custom nameservers")
	cmd.Flags().BoolVar(&ipv6, "ipv6", false, "include AAAA (IPv6) records")

	return cmd
}

// buildCLINameserverList constructs the final list of nameservers for CLI usage.
func buildCLINameserverList(nameservers []string, includeDefault bool) []string {
	defaults := dnslookup.GetDefaultNameservers()

	if len(nameservers) == 0 {
		return defaults
	}

	if includeDefault {
		return append(defaults, nameservers...)
	}

	return nameservers
}

// Execute adds all child commands to the root command and sets flags appropriately.
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
}

// printRecords prints DNS records to the console, one table per nameserver.
func printRecords(w io.Writer, results []models.NameserverResult) {
	for i, nsResult := range results {
		if i > 0 {
			fmt.Fprintln(w)
		}

		// Print nameserver header
		nsDisplay := dnslookup.FormatNameserver(nsResult.Nameserver)
		fmt.Fprintf(w, "Nameserver: %s\n", nsDisplay)

		// Show error if nameserver is unreachable
		if nsResult.Error != "" {
			fmt.Fprintf(w, "  ✗ Error: %s\n", nsResult.Error)
			continue
		}

		tw := table.NewWriter()
		tw.SetStyle(table.StyleLight)
		tw.Style().Options.DrawBorder = true
		tw.Style().Options.SeparateHeader = true
		tw.Style().Options.SeparateRows = true
		tw.Style().Options.SeparateColumns = true

		tw.AppendHeader(table.Row{"Domain", "Host", "Type", "Data"})

		hasRows := false
		for domain, records := range nsResult.Results {
			if len(records) == 0 {
				tw.AppendRow(table.Row{domain, "", "", "No record found"})
				hasRows = true
			} else {
				for _, r := range records {
					tw.AppendRow(table.Row{domain, r.Host, r.Type, r.Data})
					hasRows = true
				}
			}
		}

		tw.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AutoMerge: true},
		})

		if hasRows {
			fmt.Fprintln(w, tw.Render())
		}
	}
}

// FormatResultsAsText formats results as plain text tables (for copy-as-text in web UI).
func FormatResultsAsText(results []models.NameserverResult) string {
	var sb strings.Builder
	printRecords(&sb, results)
	return sb.String()
}
