package dns

import (
	"errors"
	"log"
	"net"
	"strconv"

	"github.com/mfuentesg/localdns/storage"
	"github.com/miekg/dns"
)

type Option func(*handler)

func WithStorage(st storage.Storage) Option {
	return func(s *handler) {
		s.st = st
	}
}

func WithPort(port int) Option {
	return func(h *handler) {
		h.port = port
	}
}

func WithHost(host string) Option {
	return func(h *handler) {
		h.host = host
	}
}

func WithProtocol(protocol string) Option {
	return func(h *handler) {
		h.protocol = protocol
	}
}

type handler struct {
	st       storage.Storage
	port     int
	host     string
	protocol string
}

func writeMessage(w dns.ResponseWriter, message *dns.Msg) {
	_ = w.WriteMsg(message)
}

func (s *handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	var message dns.Msg
	question := m.Question[0]

	if question.Qtype != dns.TypeA {
		_ = w.WriteMsg(&message)
		return
	}

	message.SetReply(m)
	domain := question.Name
	record, err := s.st.Get(domain)

	message.Authoritative = true

	if errors.Is(err, storage.ErrRecordNotFound) {
		log.Printf("domain '%s' not found\n", domain)
		writeMessage(w, &message)
		return
	}

	if err != nil {
		log.Printf("could not process domain '%s', reason: %+v\n", domain, err)
		// TODO: check if writeMessage is required
		writeMessage(w, &message)
		return
	}

	message.Answer = append(message.Answer, &dns.A{
		Hdr: dns.RR_Header{
			Name:   domain,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    60,
		},
		A: net.ParseIP(record.IP),
	})

	writeMessage(w, &message)
}

func New(opts ...Option) dns.Server {
	h := handler{port: 53, host: "", protocol: "udp"}

	for _, opt := range opts {
		opt(&h)
	}

	return dns.Server{
		Addr:    net.JoinHostPort(h.host, strconv.Itoa(h.port)),
		Net:     h.protocol,
		Handler: &h,
	}
}
