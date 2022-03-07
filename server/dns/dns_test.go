package dns

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type fakeHandler struct{}

func (h *fakeHandler) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	_ = w.WriteMsg(m)
}

func TestNew(t *testing.T) {
	server := New(new(fakeHandler))

	assert.IsType(t, new(dns.Server), server)
	assert.Equal(t, server.Net, "udp")
	assert.Equal(t, server.Addr, ":53")
}

func TestWithAddr(t *testing.T) {
	addr := "0.0.0.0:8053"
	server := New(new(fakeHandler), WithAddr(addr))

	assert.Equal(t, server.Addr, addr)
}

func TestWithProtocol(t *testing.T) {
	protocol := "tcp"
	server := New(new(fakeHandler), WithProtocol(protocol))

	assert.Equal(t, server.Net, protocol)
}
