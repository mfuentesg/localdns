package handler

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"testing"

	"github.com/mfuentesg/dnstest"
	"github.com/mfuentesg/localdns/storage"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type fakeStorage struct {
	storage.Storage
}

func (fs *fakeStorage) Get(identifier string) (*storage.Record, error) {
	if identifier == dns.Fqdn("www.google.com") {
		return nil, errors.New("fake error")
	}

	if identifier == dns.Fqdn("www.fake.com") {
		return nil, storage.ErrRecordNotFound
	}

	return &storage.Record{
		Type:   "A",
		Domain: identifier,
		IP:     "192.168.1.10",
		TTL:    1000,
	}, nil
}

type fakeDNSRW struct {
	dns.ResponseWriter
}

func (w *fakeDNSRW) WriteMsg(_ *dns.Msg) error {
	return nil
}

func TestNew(t *testing.T) {
	h := New(new(fakeStorage))

	assert.NotNil(t, h)
	assert.IsType(t, new(Handler), h)
}

func TestWithDNSServer(t *testing.T) {
	dnsServer := "1.1.1.1:53"
	h := New(new(fakeStorage), WithDNSServer(dnsServer))

	assert.Equal(t, dnsServer, h.dnsServer)
}

func TestHandler_ServeDNS(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("unsupported dns type", func(tt *testing.T) {
		h := New(new(fakeStorage))
		message := new(dns.Msg)
		message.SetQuestion(dns.Fqdn("www.google.com"), dns.TypeAAAA)

		recorder := dnstest.NewRecorder(new(fakeDNSRW))
		h.ServeDNS(recorder, message)
		assert.Len(tt, recorder.Msg.Answer, 0)
	})

	t.Run("unexpected storage error", func(tt *testing.T) {
		h := New(new(fakeStorage))
		message := new(dns.Msg)
		message.SetQuestion(dns.Fqdn("www.google.com"), dns.TypeA)

		recorder := dnstest.NewRecorder(new(fakeDNSRW))
		h.ServeDNS(recorder, message)
		assert.Len(tt, recorder.Msg.Answer, 0)
	})

	t.Run("forward query to DNS server", func(tt *testing.T) {
		message := new(dns.Msg)
		message.SetQuestion(dns.Fqdn("www.fake.com"), dns.TypeA)

		server := dnstest.NewServer(func(writer dns.ResponseWriter, msg *dns.Msg) {
			m := new(dns.Msg)
			m.SetReply(msg)

			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   "www.fake.com.",
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				A: net.ParseIP("1.1.1.1"),
			})

			_ = writer.WriteMsg(m)
		})
		defer server.Close()
		h := New(new(fakeStorage), WithDNSServer(server.Addr))
		recorder := dnstest.NewRecorder(new(fakeDNSRW))

		h.ServeDNS(recorder, message)
		assert.Len(tt, recorder.Msg.Answer, 1)
		assert.Equal(tt, "1.1.1.1", recorder.Msg.Answer[0].(*dns.A).A.String())
	})

	t.Run("get stored dns from storage", func(tt *testing.T) {
		message := new(dns.Msg)
		message.SetQuestion(dns.Fqdn("www.valid-domain.com"), dns.TypeA)

		h := New(new(fakeStorage))
		recorder := dnstest.NewRecorder(new(fakeDNSRW))

		h.ServeDNS(recorder, message)
		assert.Len(tt, recorder.Msg.Answer, 1)
		assert.Equal(tt, "192.168.1.10", recorder.Msg.Answer[0].(*dns.A).A.String())
	})
}

func BenchmarkHandler_ServeDNS(b *testing.B) {}
