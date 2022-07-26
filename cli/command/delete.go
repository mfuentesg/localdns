package command

import (
	"context"

	"github.com/mfuentesg/localdns/cli/transport"
	"github.com/mfuentesg/localdns/pb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete given record by id",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("At least one record ID must be provided")
		}

		server, _ := cmd.Flags().GetString("server")
		client, err := transport.New(server)
		if err != nil {
			log.Fatal(err)
		}

		for _, id := range args {
			if _, err := client.DeleteRecord(context.Background(), &pb.Record{Id: id}); err != nil {
				log.Error(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
