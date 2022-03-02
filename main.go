package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mfuentesg/localdns/server/dns"
	"github.com/mfuentesg/localdns/storage/embedded"
)

func main() {
	st, err := embedded.New()
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	errs := make(chan error, 2)

	go func() {
		s := dns.New(
			dns.WithProtocol("udp"),
			dns.WithStorage(st),
			dns.WithPort(8053),
			dns.WithDNSServer("8.8.8.8:53"),
		)
		log.Printf("udp: dns server started at %s\n", s.Addr)
		errs <- s.ListenAndServe()
	}()

	go func() {
		s := dns.New(
			dns.WithProtocol("tcp"),
			dns.WithStorage(st),
			dns.WithPort(8053),
			dns.WithDNSServer("8.8.8.8:53"),
		)
		log.Printf("tcp: dns server started at %s\n", s.Addr)
		errs <- s.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	log.Printf("service terminated: %s\n", <-errs)
}
