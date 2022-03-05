package storage

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("localdns: record not found")
)

type Record struct {
	Type   string
	Domain string
	IP     string
	TTL    time.Duration
}

type Storage interface {
	Put(r Record) error
	Get(identifier string) (*Record, error)
	Delete(identifier string) error
	Close() error
}
