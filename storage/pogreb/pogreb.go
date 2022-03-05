package pogreb

import (
	"bytes"
	"encoding/gob"

	"github.com/akrylysov/pogreb"
	"github.com/mfuentesg/localdns/storage"
)

type Pogreb struct {
	db *pogreb.DB
}

func New() (*Pogreb, error) {
	db, err := pogreb.Open("dns", nil)
	if err != nil {
		return nil, err
	}

	return &Pogreb{db: db}, nil
}

func (pg *Pogreb) Put(r storage.Record) error {
	var record bytes.Buffer
	if err := gob.NewEncoder(&record).Encode(r); err != nil {
		return err
	}

	if err := pg.db.Put([]byte(r.Domain), record.Bytes()); err != nil {
		return err
	}

	return nil
}

func (pg *Pogreb) Get(key string) (*storage.Record, error) {
	data, err := pg.db.Get([]byte(key))

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

func (pg *Pogreb) Delete(key string) error {
	return pg.db.Delete([]byte(key))
}

func (pg *Pogreb) Close() error {
	return pg.db.Close()
}
