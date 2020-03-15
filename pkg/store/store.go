package store

import (
	"net"
	"os"
	"sync"
	"time"

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
func (s *SqliteStore) PutData(addr net.Addr, packetSize int, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("INSERT INTO data (timestamp, from_addr, packet_size, payload) VALUES(?,?,?,?)", time.Now(), addr.String(), packetSize, data)
	return err
}
