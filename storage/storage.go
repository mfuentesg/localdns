package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/akrylysov/pogreb"
)

type Record struct {
	Type   string
	Domain string
	IP     string
	TTL    time.Duration
}

type Storage struct {
	db *pogreb.DB
}

func New() (*Storage, error) {
	db, err := pogreb.Open("dns", nil)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (d *Storage) Close() {
	defer func() { _ = d.db.Close() }()
}

func (d *Storage) Put(r Record) error {
	var record bytes.Buffer
	if err := gob.NewEncoder(&record).Encode(r); err != nil {
		return fmt.Errorf("could not encode the given record, %+v", err)
	}

	if err := d.db.Put([]byte(r.Domain), record.Bytes()); err != nil {
		return fmt.Errorf("could not put the given record, %+v", err)
	}

	return nil
}

func (d *Storage) Get(key string) (*Record, error) {
	data, err := d.db.Get([]byte(key))

	if err != nil {
		return nil, fmt.Errorf("could not retrieve the record %+v", err)
	}

	var record Record
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&record); err != nil {
		return nil, fmt.Errorf("could not decode the given record, %+v", err)
	}

	return &record, nil
}

func (d *Storage) Delete(key string) error {
	if err := d.db.Delete([]byte(key)); err != nil {
		return fmt.Errorf("could not delete the record, %+v", err)
	}

	return nil
}
