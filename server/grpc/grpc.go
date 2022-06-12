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
		Ipv4:   r.Ipv4,
		Ipv6:   r.Ipv6,
		Ttl:    r.Ttl,
	}

	if !strings.HasSuffix(record.Domain, ".") {
		record.Domain += "."
	}

	if record.Ttl == 0 {
		record.Ttl = 604800
	}

	id, err := srv.st.Put(storage.Record{
		Type:   record.Type,
		Domain: record.Domain,
		IPv4:   record.Ipv4,
		IPv6:   record.Ipv6,
		TTL:    record.Ttl,
	})

	logEntry := log.WithFields(log.Fields{
		"domain": record.Domain,
		"id":     id,
		"type":   record.Type,
		"ipv4":   record.Ipv4,
		"ipv6":   record.Ipv6,
	})

	if err != nil {
		logEntry.WithField("reason", err).Error("unable to store data")
		return nil, err
	}

	record.Id = id

	logEntry.Info("record created")
	return srv.GetRecord(ctx, &record)
}

func (srv *Server) DeleteRecord(_ context.Context, r *pb.Record) (*emptypb.Empty, error) {
	err := srv.st.Delete(r.Id)
	logEntry := log.WithFields(log.Fields{"id": r.Id})
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
			Id:     record.ID,
			Type:   record.Type,
			Domain: record.Domain,
			Ipv4:   record.IPv4,
			Ipv6:   record.IPv6,
			Ttl:    record.TTL,
		})
	}

	return &pb.RecordList{Records: list}, nil
}

func (srv *Server) GetRecord(_ context.Context, r *pb.Record) (*pb.Record, error) {
	record, err := srv.st.Get(r.Id)
	logEntry := log.WithFields(log.Fields{
		"id": r.Id,
	})

	if err != nil {
		logEntry.Error("unable to retrieve the record")
		return nil, err
	}

	logEntry.Info("record retrieved")
	return &pb.Record{
		Id:     record.ID,
		Type:   record.Type,
		Domain: record.Domain,
		Ipv4:   record.IPv4,
		Ipv6:   record.IPv6,
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
