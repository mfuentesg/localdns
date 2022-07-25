package transport

import (
	"github.com/mfuentesg/localdns/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb.DnsServiceClient
	conn *grpc.ClientConn
}

func New(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:             conn,
		DnsServiceClient: pb.NewDnsServiceClient(conn),
	}, nil
}

func (cli *Client) Close() error {
	return cli.conn.Close()
}
