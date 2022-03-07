package database

import (
	"testing"

	"github.com/mfuentesg/localdns/storage/pogreb"
	"github.com/mfuentesg/localdns/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("unknown engine", func(tt *testing.T) {
		st, err := New(10)

		assert.Error(tt, err)
		assert.Nil(tt, st)
	})

	t.Run("return expected pogreb engine", func(tt *testing.T) {
		st, err := New(PogrebEngine)

		assert.Nil(tt, err)
		assert.NotNil(tt, st)
		assert.IsType(tt, new(pogreb.Pogreb), st)
	})

	t.Run("return expected sqlite engine", func(tt *testing.T) {
		st, err := New(SQLiteEngine)

		assert.Nil(tt, err)
		assert.NotNil(tt, st)
		assert.IsType(tt, new(sqlite.SQLite), st)
	})
}
