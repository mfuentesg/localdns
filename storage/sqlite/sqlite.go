package sqlite

import (
	"errors"

	"github.com/mfuentesg/localdns/storage"
)

type SQLite struct{}

func New() (*SQLite, error) {
	return nil, errors.New("unimplemented")
}

func (sq *SQLite) Put(r storage.Record) error {
	return errors.New("unimplemented")
}

func (sq *SQLite) Get(identifier string) (*storage.Record, error) {
	return nil, errors.New("unimplemented")
}

func (sq *SQLite) Delete(identifier string) error {
	return errors.New("unimplemented")
}

func (sq *SQLite) Close() error {
	return errors.New("unimplemented")
}
