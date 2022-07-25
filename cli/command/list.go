package command

import (
	"context"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mfuentesg/localdns/cli/transport"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all registered dns records",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := cmd.Flags().GetString("server")
		if err != nil {
			log.Fatal(err)
		}

		client, err := transport.New(server)
		if err != nil {
			log.Fatal(err)
		}

		records, err := client.ListRecords(context.Background(), new(emptypb.Empty))
		if err != nil {
			log.Fatal(err)
		}

		if len(records.Records) == 0 {
			log.Info("No records found")
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
}
