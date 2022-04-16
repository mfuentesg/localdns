package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mfuentesg/localdns/handler"
	"github.com/mfuentesg/localdns/server/dns"
	"github.com/mfuentesg/localdns/server/grpc"
	"github.com/mfuentesg/localdns/storage/sqlite"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetLevel(log.InfoLevel)

	db, err := sqlite.New()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	errs := make(chan error, 2)

	go func() {
		protocol := "udp"
		h := handler.New(db, handler.WithDNSServer("8.8.8.8:53"), handler.WithProtocol(protocol))
		s := dns.New(h, dns.WithAddr(":8053"), dns.WithProtocol(protocol))
		log.WithFields(log.Fields{"protocol": protocol, "addr": s.Addr}).Info("dns server started")
		errs <- s.ListenAndServe()
	}()

	go func() {
		protocol := "tcp"
		h := handler.New(db, handler.WithDNSServer("8.8.8.8:53"), handler.WithProtocol(protocol))
		s := dns.New(h, dns.WithAddr(":8053"), dns.WithProtocol(protocol))
		log.WithFields(log.Fields{"protocol": protocol, "addr": s.Addr}).Info("dns server started")
		errs <- s.ListenAndServe()
	}()

	go func() {
		s := grpc.New(db, grpc.WithAddr(":8080"))
		log.WithField("addr", s.Addr).Info("grpc server started")
		errs <- s.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	log.Errorf("localdns: service %s", <-errs)
}
