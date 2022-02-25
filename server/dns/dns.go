package dns

import (
	"fmt"
	"net"

	"github.com/mfuentesg/localdns/storage"
	"github.com/miekg/dns"
)

type Option func(*handler)

func WithStorage(st *storage.Storage) Option {
	return func(s *handler) {
		s.st = st
	}
}

type handler struct {
	st *storage.Storage
}

func (s *handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	var message dns.Msg

	message.SetReply(m)
	question := m.Question[0]
	domain := question.Name

	if question.Qtype == dns.TypeA {
		record, err := s.st.Get(question.Name)

		message.Authoritative = true
		message.RecursionAvailable = true

		if err == nil {
			message.Answer = append(message.Answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   domain,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				A: net.ParseIP(record.IP),
			})
		} else {
			fmt.Printf("could not handle record %+v\n", err)
		}
	}

	_ = w.WriteMsg(&message)
}

func New(opts ...Option) dns.Server {
	var h handler
	for _, opt := range opts {
		opt(&h)
	}

	return dns.Server{
		Addr:    ":53",
		Net:     "udp",
		Handler: &h,
	}
}
