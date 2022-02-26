package storage

import (
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("localdns: record not found")

type Record struct {
	Type   string
	Domain string
	IP     string
	TTL    time.Duration
}

type Storage interface {
	Put(r Record) error
	Get(key string) (*Record, error)
	Delete(key string) error
}
