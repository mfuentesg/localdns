package dns

import (
	"github.com/miekg/dns"
)

type Server struct {
	Addr     string
	Protocol string
}

type Option func(*Server)

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

func WithProtocol(protocol string) Option {
	return func(s *Server) {
		s.Protocol = protocol
	}
}

func New(handler dns.Handler, opts ...Option) *dns.Server {
	srv := Server{
		Addr:     ":53",
		Protocol: "udp",
	}

	for _, opt := range opts {
		opt(&srv)
	}

	return &dns.Server{
		Addr:    srv.Addr,
		Net:     srv.Protocol,
		Handler: handler,
	}
}
