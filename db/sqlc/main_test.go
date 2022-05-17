package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driver = "postgres"
	source = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(driver, source)
	if err != nil {
		log.Fatal("can't connect to database:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}