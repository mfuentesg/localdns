package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mfuentesg/localdns/handler"
	"github.com/mfuentesg/localdns/server/dns"
	"github.com/mfuentesg/localdns/server/grpc"
	"github.com/mfuentesg/localdns/storage/sqlite"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	db, err := sqlite.New()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	h := handler.New(db, handler.WithDNSServer("8.8.8.8:53"))
	errs := make(chan error, 2)

	go func() {
		s := dns.New(h, dns.WithAddr(":8053"), dns.WithProtocol("udp"))
		log.Printf("udp: dns server started at %s\n", s.Addr)
		errs <- s.ListenAndServe()
	}()

	go func() {
		s := dns.New(h, dns.WithAddr(":8053"), dns.WithProtocol("tcp"))
		log.Printf("tcp: dns server started at %s\n", s.Addr)
		errs <- s.ListenAndServe()
	}()

	go func() {
		s := grpc.New(grpc.WithStorage(db), grpc.WithAddr(":8080"))
		log.Printf("grpc server started at %s\n", s.Addr)
		errs <- s.ListenAndServe()
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("prometheus server started at :9090")
		errs <- http.ListenAndServe(":9090", nil)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	log.Printf("localdns: service %s\n", <-errs)
}
