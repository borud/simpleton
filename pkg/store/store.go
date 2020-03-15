package store

import (
	"os"
	"sync"

	"github.com/borud/simpleton/pkg/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
)

// SqliteStore implements the store interface with Sqlite
type SqliteStore struct {
	mu sync.Mutex
	db *sqlx.DB
}

// New creates new Store backed by SQLite3
func New(dbFile string) (*SqliteStore, error) {
	var databaseFileExisted = false
	if _, err := os.Stat(dbFile); err == nil {
		databaseFileExisted = true
	}

	d, err := sqlx.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if err = d.Ping(); err != nil {
		return nil, err
	}

	if !databaseFileExisted {
		createSchema(d, dbFile)
	}

	return &SqliteStore{db: d}, nil
}

// PutData stores one row of data in the database
func (s *SqliteStore) PutData(data *model.Data) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, err := s.db.NamedExec("INSERT INTO data (timestamp, from_addr, packet_size, payload) VALUES(:timestamp,:from_addr,:packet_size,:payload)", data)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

// ListData returns a list of the data from the database sorted by ID
// in descending order (newest first)
func (s *SqliteStore) ListData(offset int, limit int) ([]model.Data, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var data []model.Data
	err := s.db.Select(&data, "SELECT * FROM data ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
	return data, err
}

// Get fetches a single datapoint by id
func (s *SqliteStore) Get(id int64) (*model.Data, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var row model.Data
	err := s.db.QueryRowx("SELECT * FROM data WHERE id = ?", id).StructScan(&row)
	return &row, err
}
