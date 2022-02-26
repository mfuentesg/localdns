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

	s := dns.New(dns.WithStorage(st))
	errs := make(chan error, 2)

	go func() {
		log.Printf("dns server started\n")
		errs <- s.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	log.Printf("service terminated: %s\n", <-errs)
}
