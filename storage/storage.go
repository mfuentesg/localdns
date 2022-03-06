package storage

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("localdns: record not found")
)

type Record struct {
	Type   string
	Domain string
	IP     string
	TTL    int32
}

type Storage interface {
	Put(r Record) error
	Get(identifier string) (*Record, error)
	Delete(identifier string) error
	List() ([]*Record, error)
	Close() error
}
