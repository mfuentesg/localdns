package grpc

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/mfuentesg/localdns/pb"
	"github.com/mfuentesg/localdns/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Option func(*Server)

type Server struct {
	pb.DnsServiceServer
	st   storage.Storage
	Addr string
}

func WithStorage(st storage.Storage) Option {
	return func(s *Server) {
		s.st = st
	}
}

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

func (srv *Server) PutRecord(ctx context.Context, r *pb.Record) (*pb.Record, error) {
	domain := r.Domain
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	ttl := r.Ttl
	if ttl == 0 {
		ttl = 604800
	}

	err := srv.st.Put(storage.Record{
		Type:   r.Type,
		Domain: domain,
		IP:     r.Ip,
		TTL:    ttl,
	})

	if err != nil {
		return nil, err
	}

	return srv.GetRecord(ctx, r)
}

func (srv *Server) DeleteRecord(ctx context.Context, r *pb.Record) (*emptypb.Empty, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (srv *Server) ListRecords(ctx context.Context, f *pb.RecordsFilter) (*pb.RecordList, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (srv *Server) GetRecord(_ context.Context, r *pb.Record) (*pb.Record, error) {
	record, err := srv.st.Get(r.Domain)
	if err != nil {
		return nil, err
	}

	return &pb.Record{
		Type:   record.Type,
		Domain: record.Domain,
		Ip:     record.IP,
		Ttl:    record.TTL,
	}, nil
}

func (srv *Server) ListenAndServe() error {
	server := grpc.NewServer()

	reflection.Register(server)
	pb.RegisterDnsServiceServer(server, srv)

	lis, err := net.Listen("tcp", srv.Addr)

	if err != nil {
		return err
	}

	return server.Serve(lis)
}

func New(opts ...Option) *Server {
	var srv Server

	for _, opt := range opts {
		opt(&srv)
	}

	return &srv
}
