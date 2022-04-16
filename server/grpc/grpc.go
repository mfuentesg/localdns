package grpc

import (
	"context"
	"net"
	"strings"

	"github.com/mfuentesg/localdns/pb"
	"github.com/mfuentesg/localdns/storage"
	log "github.com/sirupsen/logrus"
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

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

func (srv *Server) PutRecord(ctx context.Context, r *pb.Record) (*pb.Record, error) {
	record := pb.Record{
		Type:   r.Type,
		Domain: r.Domain,
		Ip:     r.Ip,
		Ttl:    r.Ttl,
	}

	if !strings.HasSuffix(record.Domain, ".") {
		record.Domain += "."
	}

	if record.Ttl == 0 {
		record.Ttl = 604800
	}

	err := srv.st.Put(storage.Record{
		Type:   record.Type,
		Domain: record.Domain,
		IP:     record.Ip,
		TTL:    record.Ttl,
	})

	logEntry := log.WithFields(log.Fields{
		"domain": record.Domain,
		"type":   record.Type,
		"ip":     record.Ip,
	})

	if err != nil {
		logEntry.Error("unable to store data")
		return nil, err
	}

	logEntry.Info("record created")
	return srv.GetRecord(ctx, &record)
}

func (srv *Server) DeleteRecord(_ context.Context, r *pb.Record) (*emptypb.Empty, error) {
	err := srv.st.Delete(r.Domain)
	logEntry := log.WithFields(log.Fields{
		"domain": r.Domain,
	})
	if err != nil {
		logEntry.Error("unable to delete the record")
		return nil, err
	}
	return new(emptypb.Empty), nil
}

func (srv *Server) ListRecords(_ context.Context, _ *emptypb.Empty) (*pb.RecordList, error) {
	records, err := srv.st.List()

	if err != nil {
		log.Error("unable to retrieve records")
		return nil, err
	}

	list := make([]*pb.Record, 0, len(records))
	for _, record := range records {
		list = append(list, &pb.Record{
			Type:   record.Type,
			Domain: record.Domain,
			Ip:     record.IP,
			Ttl:    record.TTL,
		})
	}

	return &pb.RecordList{Records: list}, nil
}

func (srv *Server) GetRecord(_ context.Context, r *pb.Record) (*pb.Record, error) {
	record, err := srv.st.Get(r.Domain)
	logEntry := log.WithFields(log.Fields{
		"domain": record.Domain,
		"type":   record.Type,
		"ip":     record.IP,
	})

	if err != nil {
		logEntry.Error("unable to retrieve the record")
		return nil, err
	}

	logEntry.Info("record retrieved")
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

func New(db storage.Storage, opts ...Option) *Server {
	srv := &Server{Addr: ":8080", st: db}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}
