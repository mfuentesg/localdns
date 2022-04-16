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

func (h *Handler) forwardQuery(message *dns.Msg) (*dns.Msg, error) {
	c := &dns.Client{Net: h.protocol}
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
		"dnsType":  question.Qtype,
		"domain":   domain,
		"question": question,
		"protocol": h.protocol,
	})

	if question.Qtype != dns.TypeA {
		logEntry.Info("unsupported dns type")
		return nil, errors.New("unsupported dns type")
	}

	record, err := h.storage.Get(domain)

	if errors.Is(err, storage.ErrRecordNotFound) {
		logEntry.Info("unregistered domain")

		forwardedMessage, err := h.forwardQuery(m)
		if err != nil {
			logEntry.WithField("reason", err).Error("unable to forward query")
		}

		logEntry.Info("record forwarded successfully")
		return forwardedMessage, err
	}

	if err != nil {
		logEntry.WithField("reason", err).Error("unable to get record from database")
		return nil, err
	}

	logEntry.Info("domain found")

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
