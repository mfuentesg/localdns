package main

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	sq "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/mfuentesg/localdns/handler"
	"github.com/mfuentesg/localdns/server/dns"
	"github.com/mfuentesg/localdns/server/grpc"
	"github.com/mfuentesg/localdns/storage/sqlite"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed migrations
var migrations embed.FS

func loadConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName(".localdns")
	viper.SetConfigType("yaml")

	viper.SetDefault("remote_server", "8.8.8.8:53")
	viper.SetDefault("servers.dns_udp.addr", ":8053")
	viper.SetDefault("servers.dns_udp.enabled", true)
	viper.SetDefault("servers.dns_tcp.addr", ":8053")
	viper.SetDefault("servers.dns_tcp.enabled", true)
	viper.SetDefault("servers.grpc.addr", ":8080")
	viper.SetDefault("servers.grpc.enabled", true)
	viper.SetDefault("servers.prometheus.addr", ":9090")
	viper.SetDefault("servers.prometheus.enabled", true)
	viper.SetDefault("database.dsn", "localdns.db")

	_ = viper.ReadInConfig()
	_ = viper.MergeInConfig()

	return viper.WriteConfigAs(".localdns.yaml")
}

func applyMigrations() error {
	db, err := sql.Open("sqlite", viper.GetString("database.dsn"))
	if err != nil {
		return err
	}
	sourceDriver, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return err
	}
	target, err := sq.WithInstance(db, new(sq.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance(
		"httpfs", sourceDriver, "sqlite", target)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	defer db.Close()
	return sourceDriver.Close()
}

func main() {
	log.SetFormatter(new(log.JSONFormatter))
	log.SetLevel(log.InfoLevel)

	if err := loadConfig(); err != nil {
		log.WithField("reason", err).Fatal("unable to read the config file")
	}

	if err := applyMigrations(); err != nil {
		log.WithField("reason", err).Fatal("unable to apply migrations")
	}

	db, err := sqlite.New(viper.GetString("database.dsn"))
	if err != nil {
		log.WithField("reason", err).Fatal("unable to load db")
	}

	defer db.Close()

	errs := make(chan error, 2)

	udpEnabled := viper.GetBool("servers.dns_udp.enabled")
	tcpEnabled := viper.GetBool("servers.dns_tcp.enabled")
	grpcEnabled := viper.GetBool("servers.grpc.enabled")

	if udpEnabled {
		go func() {
			protocol := "udp"
			h := handler.New(db, handler.WithDNSServer(viper.GetString("remote_server")), handler.WithProtocol(protocol))
			s := dns.New(h, dns.WithAddr(viper.GetString("servers.dns_udp.addr")), dns.WithProtocol(protocol))
			log.WithFields(log.Fields{"protocol": protocol, "addr": s.Addr}).Info("dns server started")
			errs <- s.ListenAndServe()
		}()
	}

	if tcpEnabled {
		go func() {
			protocol := "tcp"
			h := handler.New(db, handler.WithDNSServer(viper.GetString("remote_server")), handler.WithProtocol(protocol))
			s := dns.New(h, dns.WithAddr(viper.GetString("servers.dns_tcp.addr")), dns.WithProtocol(protocol))
			log.WithFields(log.Fields{"protocol": protocol, "addr": s.Addr}).Info("dns server started")
			errs <- s.ListenAndServe()
		}()
	}

	if grpcEnabled {
		go func() {
			s := grpc.New(db, grpc.WithAddr(viper.GetString("servers.grpc.addr")))
			log.WithField("addr", s.Addr).Info("grpc server started")
			errs <- s.ListenAndServe()
		}()
	}

	if grpcEnabled || udpEnabled || tcpEnabled {
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			errs <- fmt.Errorf("%s", <-c)
		}()
		log.Errorf("localdns: service %s", <-errs)
	} else {
		log.Info("there are not enabled services in config file")
	}
}
