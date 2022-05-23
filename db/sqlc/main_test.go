package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ahmadfarhanstwn/simple_bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	testDB, err = sql.Open(config.DB_Driver, config.DB_Source)
	if err != nil {
		log.Fatal("can't connect to database:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}