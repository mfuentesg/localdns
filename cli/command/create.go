package command

import (
	"context"

	"github.com/mfuentesg/localdns/cli/transport"
	"github.com/mfuentesg/localdns/pb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new record",
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		client, err := transport.New(server)
		if err != nil {
			log.Fatal(err)
		}

		recordType, _ := cmd.Flags().GetString("type")
		domain, _ := cmd.Flags().GetString("domain")
		ipv4, _ := cmd.Flags().GetString("ipv4")
		ipv6, _ := cmd.Flags().GetString("ipv6")
		ttl, _ := cmd.Flags().GetInt32("ttl")

		record, err := client.PutRecord(context.Background(), &pb.Record{
			Type:   recordType,
			Domain: domain,
			Ipv4:   ipv4,
			Ttl:    ttl,
			Ipv6:   ipv6,
		})

		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Record `%s` created\n", record.Id)
	},
}

func init() {
	createCmd.Flags().StringP("domain", "d", "", "record domain")
	createCmd.Flags().StringP("type", "t", "A", "record type [A|AAAA]")
	createCmd.Flags().StringP("ipv4", "4", "", "record ip version 4")
	createCmd.Flags().StringP("ipv6", "6", "", "record ip version 6")
	createCmd.Flags().Int32P("ttl", "l", 604800, "record ttl")

	_ = createCmd.MarkFlagRequired("domain")
	_ = createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagsMutuallyExclusive("ipv4", "ipv6")

	rootCmd.AddCommand(createCmd)
}
