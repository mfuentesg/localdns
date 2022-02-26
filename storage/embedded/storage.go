package embedded

import (
	"bytes"
	"encoding/gob"

	"github.com/akrylysov/pogreb"
	"github.com/mfuentesg/localdns/storage"
)

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

func (d *Storage) Put(r storage.Record) error {
	var record bytes.Buffer
	if err := gob.NewEncoder(&record).Encode(r); err != nil {
		return err
	}

	if err := d.db.Put([]byte(r.Domain), record.Bytes()); err != nil {
		return err
	}

	return nil
}

func (d *Storage) Get(key string) (*storage.Record, error) {
	data, err := d.db.Get([]byte(key))

	if len(data) == 0 {
		return nil, storage.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	var record storage.Record
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (d *Storage) Delete(key string) error {
	return d.db.Delete([]byte(key))
}
