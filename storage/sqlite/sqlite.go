package sqlite

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/mfuentesg/localdns/storage"
	_ "modernc.org/sqlite" // sqlite driver
)

type SQLite struct {
	db *sqlx.DB
}

func (sq *SQLite) prepareDB() error {
	query := `
		create table if not exists records
		(
			id         varchar(36) default (lower(hex(randomblob(4)) || '-' || hex(randomblob(2))
				|| '-' || '4' || substr(hex(randomblob(2)), 2) || '-'
				|| substr('AB89', 1 + (abs(random()) % 4), 1) ||
												  substr(hex(randomblob(2)), 2) || '-' || hex(randomblob(6)))) primary key,
			domain     text,
			ipv4       varchar(15),
			ipv6       varchar(39),
			created_at datetime    default CURRENT_TIMESTAMP,
			ttl        integer     default 604800,
			type       varchar(10)
		);
	`

	_, err := sq.db.Exec(query)
	return err
}

func New() (*SQLite, error) {
	db, err := sqlx.Open("sqlite", "localdns.db")
	if err != nil {
		return nil, err
	}

	sq := &SQLite{db: db}

	if err := sq.prepareDB(); err != nil {
		return nil, err
	}

	return sq, nil
}

func (sq *SQLite) Put(r storage.Record) (string, error) {
	query := `insert or replace into records(
    	domain, ipv4, ipv6, ttl, type
	) values(?, ?, ?, ?, ?) returning id`

	var id string
	err := sq.db.QueryRow(query, r.Domain, r.IPv4, r.IPv6, r.TTL, r.Type).Scan(&id)
	return id, err
}

func (sq *SQLite) Get(identifier string) (*storage.Record, error) {
	var record storage.Record
	query := `select id, domain, ipv4, ipv6, ttl, type from records where id = ?`
	err := sq.db.QueryRow(query, identifier).
		Scan(&record.ID, &record.Domain, &record.IPv4, &record.IPv6, &record.TTL, &record.Type)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (sq *SQLite) Delete(identifier string) error {
	query := `delete from records where id = ?`
	_, err := sq.db.Exec(query, identifier)

	return err
}

func (sq *SQLite) List() ([]*storage.Record, error) {
	var records []*storage.Record
	query := `select id, domain, ipv4, ipv6, ttl, type from records`
	if err := sq.db.Select(&records, query); err != nil {
		return nil, err
	}
	return records, nil
}

func (sq *SQLite) Close() error {
	return sq.db.Close()
}
