package store

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Load sqlite3 driver
)

const schema = `
CREATE TABLE IF NOT EXISTS data (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp    DATETIME NOT NULL,
    from_addr    STRING NOT NULL,
    packet_size  INT NOT NULL,
    payload      BLOB NOT NULL
);
`

func createSchema(db *sqlx.DB, fileName string) {
	log.Printf("Creating database schema in %s", fileName)

	for n, statement := range strings.Split(schema, ";") {
		if _, err := db.Exec(statement); err != nil {
			panic(fmt.Sprintf("Statement %d failed: \"%s\" : %s", n+1, statement, err))
		}
	}
}
