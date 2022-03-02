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

func WithDNSServer(dnsServer string) Option {
	return func(h *handler) {
		h.dnsServer = dnsServer
	}
}

type handler struct {
	st        storage.Storage
	port      int
	host      string
	protocol  string
	dnsServer string
}

func forwardQuery(message *dns.Msg) (*dns.Msg, error) {
	c := &dns.Client{Net: "udp"}
	conn, err := c.Dial("8.8.8.8:53")

	if err != nil {
		return nil, err
	}

	// TODO: should dial connection be closed?
	// defer func() { _ = conn.Close() }()

	response, _, err := c.ExchangeWithConn(message, conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (h *handler) buildMessage(m *dns.Msg) (*dns.Msg, error) {
	question := m.Question[0]

	if question.Qtype != dns.TypeA {
		return nil, errors.New("unsupported dns type")
	}

	domain := question.Name
	record, err := h.st.Get(domain)

	if errors.Is(err, storage.ErrRecordNotFound) {
		log.Printf("domain '%s' not found\n", domain)
		return forwardQuery(m)
	}

	if err != nil {
		return nil, err
	}

	var message dns.Msg
	message.SetReply(m)
	message.Authoritative = true
	message.Answer = append(message.Answer, &dns.A{
		Hdr: dns.RR_Header{
			Name:   domain,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    60,
		},
		A: net.ParseIP(record.IP),
	})

	return &message, nil
}

func (h *handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	message, err := h.buildMessage(m)

	if err != nil {
		log.Printf("could not build message, reason: %+v\n", err)
		_ = w.WriteMsg(m)
		return
	}

	_ = w.WriteMsg(message)
}

func New(opts ...Option) dns.Server {
	h := handler{
		port:      53,
		host:      "",
		protocol:  "udp",
		dnsServer: "8.8.8.8:53",
	}

	for _, opt := range opts {
		opt(&h)
	}

	return dns.Server{
		Addr:    net.JoinHostPort(h.host, strconv.Itoa(h.port)),
		Net:     h.protocol,
		Handler: &h,
	}
}
