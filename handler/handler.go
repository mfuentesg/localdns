package handler

import (
	"errors"
	"net"

	"github.com/mfuentesg/localdns/storage"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

var (
	allowedQTypes = map[uint16]string{
		dns.TypeA:    "A",
		dns.TypeAAAA: "AAAA",
	}

	qTypesMapping = map[string]uint16{
		"A":    dns.TypeA,
		"AAAA": dns.TypeAAAA,
	}
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
	qType, ok := qTypesMapping[record.Type]
	if !ok {
		return nil
	}

	if qType == dns.TypeA {
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

	return nil
}

func buildMessage(m *dns.Msg, records []*storage.Record) *dns.Msg {
	var message dns.Msg
	message.Authoritative = true
	message.SetReply(m)

	for _, record := range records {
		message.Answer = append(message.Answer, buildRecord(record))
	}

	return &message
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

func (h *Handler) handleQuery(m *dns.Msg) (*dns.Msg, error) {
	question := m.Question[0]
	domain := question.Name
	logEntry := log.WithFields(log.Fields{
		"dnsType":   question.Qtype,
		"domain":    domain,
		"question":  question,
		"protocol":  h.protocol,
		"dnsServer": h.dnsServer,
	})

	if _, ok := allowedQTypes[question.Qtype]; !ok {
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
			return nil, err
		}

		logEntry.Info("record forwarded")
		return forwardedMessage, nil
	}

	logEntry.Info("domain found")

	return buildMessage(m, records), nil
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	message, err := h.handleQuery(m)
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
