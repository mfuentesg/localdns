package handler

import (
	"errors"
	"log"
	"net"

	"github.com/mfuentesg/localdns/storage"
	"github.com/miekg/dns"
)

type Handler struct {
	storage   storage.Storage
	dnsServer string
}

type Option func(*Handler)

func WithDNSServer(dnsServer string) Option {
	return func(h *Handler) {
		h.dnsServer = dnsServer
	}
}

func (h *Handler) forwardQuery(message *dns.Msg) (*dns.Msg, error) {
	c := &dns.Client{Net: "udp"}
	conn, err := c.Dial(h.dnsServer)

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

func (h *Handler) buildMessage(m *dns.Msg) (*dns.Msg, error) {
	question := m.Question[0]

	if question.Qtype != dns.TypeA {
		return nil, errors.New("unsupported dns type")
	}

	domain := question.Name
	record, err := h.storage.Get(domain)

	if errors.Is(err, storage.ErrRecordNotFound) {
		log.Printf("domain '%s' not found\n", domain)
		return h.forwardQuery(m)
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

func (h *Handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	message, err := h.buildMessage(m)

	if err != nil {
		log.Printf("could not build message, reason: %+v\n", err)
		_ = w.WriteMsg(m)
		return
	}

	_ = w.WriteMsg(message)
}

func New(storage storage.Storage, opts ...Option) *Handler {
	handler := Handler{
		dnsServer: "8.8.8.8:53",
		storage:   storage,
	}

	for _, opt := range opts {
		opt(&handler)
	}

	return &handler
}