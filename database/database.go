package database

import (
	"errors"

	"github.com/mfuentesg/localdns/storage"
	"github.com/mfuentesg/localdns/storage/pogreb"
	"github.com/mfuentesg/localdns/storage/sqlite"
)

var (
	ErrUnsupportedEngine = errors.New("localdns: unsupported engine")
)

const (
	PogrebEngine = iota + 1
	SQLiteEngine
)

func New(engine uint) (storage.Storage, error) {
	switch engine {
	case PogrebEngine:
		return pogreb.New()
	case SQLiteEngine:
		return sqlite.New()
	default:
		return nil, ErrUnsupportedEngine
	}
}
