package embedded

import (
	"bytes"
	"encoding/gob"

	"github.com/akrylysov/pogreb"
	"github.com/mfuentesg/localdns/storage"
)

type Embedded struct {
	db *pogreb.DB
}

func New() (*Embedded, error) {
	db, err := pogreb.Open("dns", nil)
	if err != nil {
		return nil, err
	}

	return &Embedded{db: db}, nil
}

func (st *Embedded) Close() {
	defer func() { _ = st.db.Close() }()
}

func (st *Embedded) Put(r storage.Record) error {
	var record bytes.Buffer
	if err := gob.NewEncoder(&record).Encode(r); err != nil {
		return err
	}

	if err := st.db.Put([]byte(r.Domain), record.Bytes()); err != nil {
		return err
	}

	return nil
}

func (st *Embedded) Get(key string) (*storage.Record, error) {
	data, err := st.db.Get([]byte(key))

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

func (st *Embedded) Delete(key string) error {
	return st.db.Delete([]byte(key))
}
