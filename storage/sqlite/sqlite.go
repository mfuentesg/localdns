package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3" // sqlite driver
	"github.com/mfuentesg/localdns/storage"
)

type SQLite struct {
	db *sql.DB
}

func (sq *SQLite) prepareDB() error {
	query := `create table if not exists records(
		domain     text primary key,
		ip         text,
		created_at datetime default CURRENT_TIMESTAMP,
		ttl        integer  default 604800,
		type       text
	)`

	_, err := sq.db.Exec(query)
	return err
}

func New() (*SQLite, error) {
	db, err := sql.Open("sqlite3", "localdns.db")
	if err != nil {
		return nil, err
	}

	sq := &SQLite{db: db}

	if err := sq.prepareDB(); err != nil {
		return nil, err
	}

	return sq, nil
}

func (sq *SQLite) Put(r storage.Record) error {
	query := `insert or replace into records(
    	domain, ip, ttl, type
	) values(?, ?, ?, ?)`

	_, err := sq.db.Exec(query, r.Domain, r.IP, r.TTL, r.Type)

	return err
}

func (sq *SQLite) Get(identifier string) (*storage.Record, error) {
	query := `select domain, ip, ttl, type from records where domain = ?`
	row := sq.db.QueryRow(query, identifier)

	if row.Err() == sql.ErrNoRows {
		return nil, storage.ErrRecordNotFound
	}

	var record storage.Record
	if err := row.Scan(&record.Domain, &record.IP, &record.TTL, &record.Type); err != nil {
		return nil, err
	}

	return &record, nil
}

func (sq *SQLite) Delete(identifier string) error {
	return errors.New("unimplemented")
}

func (sq *SQLite) Close() error {
	return sq.db.Close()
}
