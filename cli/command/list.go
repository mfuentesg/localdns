package command

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mfuentesg/localdns/cli/transport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/emptypb"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all registered dns records",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := transport.New(viper.GetString("server"))
		if err != nil {
			os.Exit(1)
		}

		records, err := client.ListRecords(context.Background(), new(emptypb.Empty))
		if err != nil {
			os.Exit(1)
		}

		if len(records.Records) == 0 {
			_, _ = fmt.Fprintf(os.Stdout, "No records found.\n")
			os.Exit(0)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"ID", "DOMAIN", "IPv4", "IPv6", "TYPE"})
		t.AppendSeparator()

		defer func() { _ = client.Close() }()

		for _, record := range records.Records {
			t.AppendRow([]interface{}{record.Id, record.Domain, record.Ipv4, record.Ipv6, record.Type})
		}

		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
