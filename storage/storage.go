package storage

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("localdns: record not found")
)

type Record struct {
	ID     string
	Type   string
	Domain string
	IPv4   string
	IPv6   string
	TTL    int32
}

type Storage interface {
	Put(r Record) (string, error)
	Get(identifier string) (*Record, error)
	Delete(identifier string) error
	List() ([]*Record, error)
	Close() error
}
