package handler

import (
	"errors"
	"net"

	"github.com/mfuentesg/localdns/storage"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	storage   storage.Storage
	dnsServer string
	protocol  string
}

type Option func(*Handler)

func WithDNSServer(dnsServer string) Option {
	return func(h *Handler) {
		h.dnsServer = dnsServer
	}
}

func WithProtocol(protocol string) Option {
	return func(h *Handler) {
		h.protocol = protocol
	}
}

func buildRecord(record *storage.Record) dns.RR {
	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   record.Domain,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    uint32(record.TTL),
		},
		A: net.ParseIP(record.IPv4),
	}
}

func (h *Handler) forwardQuery(message *dns.Msg) (*dns.Msg, error) {
	c := &dns.Client{Net: "udp"}
	conn, err := c.Dial(h.dnsServer)

	if err != nil {
		return nil, err
	}

	defer func() { _ = conn.Close() }()

	response, _, err := c.ExchangeWithConn(message, conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (h *Handler) buildMessage(m *dns.Msg) (*dns.Msg, error) {
	question := m.Question[0]
	domain := question.Name
	logEntry := log.WithFields(log.Fields{
		"dnsType":   question.Qtype,
		"domain":    domain,
		"question":  question,
		"protocol":  h.protocol,
		"dnsServer": h.dnsServer,
	})

	if question.Qtype != dns.TypeA {
		logEntry.Info("unsupported dns type")
		return nil, errors.New("unsupported dns type")
	}

	records, err := h.storage.GetByDomain(domain)
	if err != nil {
		logEntry.WithField("reason", err).Error("unable to get record from database")
		return nil, err
	}

	if len(records) == 0 {
		forwardedMessage, err := h.forwardQuery(m)
		if err != nil {
			logEntry.WithField("reason", err).Error("unable to forward query")
		}

		logEntry.Info("record forwarded")
		return forwardedMessage, err
	}

	logEntry.Info("domain found")

	var message dns.Msg
	message.SetReply(m)
	message.Authoritative = true

	for _, record := range records {
		message.Answer = append(message.Answer, buildRecord(record))
	}

	return &message, nil
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	message, err := h.buildMessage(m)

	if err != nil {
		_ = w.WriteMsg(m)
		return
	}

	_ = w.WriteMsg(message)
}

func New(storage storage.Storage, opts ...Option) *Handler {
	handler := Handler{
		dnsServer: "8.8.8.8:53",
		storage:   storage,
		protocol:  "udp",
	}

	for _, opt := range opts {
		opt(&handler)
	}

	return &handler
}
