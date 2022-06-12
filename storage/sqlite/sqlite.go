package sqlite

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/mfuentesg/localdns/storage"
	"modernc.org/sqlite"
)

var (
	ErrRecordAlreadyExists = errors.New("record already exists")
)

// SQLite constraint codes
// Reference: https://www.sqlite.org/rescode.html#pve

const (
	ErrCodeConstraintUnique = 2067
)

type SQLite struct {
	db *sqlx.DB
}

func New(dsn string) (*SQLite, error) {
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	return &SQLite{db: db}, nil
}

func (sq *SQLite) Put(r storage.Record) (string, error) {
	query := `insert into records(
    	domain, ipv4, ipv6, ttl, type
	) values(?, ?, ?, ?, ?) returning id`

	var id string
	err := sq.db.QueryRow(query, r.Domain, r.IPv4, r.IPv6, r.TTL, r.Type).Scan(&id)

	errCode := err.(*sqlite.Error).Code()
	if errCode == ErrCodeConstraintUnique {
		return "", ErrRecordAlreadyExists
	}
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
