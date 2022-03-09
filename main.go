package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mfuentesg/localdns/database"
	"github.com/mfuentesg/localdns/handler"
	"github.com/mfuentesg/localdns/server/dns"
	"github.com/mfuentesg/localdns/server/grpc"
)

func main() {
	st, err := database.New(database.SQLiteEngine)
	if err != nil {
		log.Fatal(err)
	}

	defer st.Close()

	h := handler.New(st, handler.WithDNSServer("8.8.8.8:53"))
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
		s := grpc.New(grpc.WithStorage(st), grpc.WithAddr(":8080"))
		log.Printf("grpc server started at %s\n", s.Addr)
		errs <- s.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	log.Printf("localdns: service %s\n", <-errs)
}
