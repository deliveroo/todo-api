package repo_test

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/deliveroo/assert-go"
	"github.com/deliveroo/todo-api/selftest/deps/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestMain(m *testing.M) {
	if flag.Parse(); testing.Short() {
		return // skip in short mode
	}

	// Connect to Postgres.
	must(postgres.Connect(), "could not connect to postgres")

	// Run tests.
	result := m.Run()

	// Reset the database.
	must(postgres.Reset(), "error resetting database")

	os.Exit(result)
}

type testDB struct {
	t    *testing.T
	pool *pgxpool.Pool
}

func getDB(t *testing.T) *testDB {
	pool, err := postgres.GetPool()
	assert.Must(t, err)
	return &testDB{
		t:    t,
		pool: pool,
	}
}

func (db *testDB) Close() {
	if stat := db.pool.Stat(); stat.IdleConns() != stat.TotalConns() {
		db.t.Error("database connection was leaked")
	}
	db.pool.Close()
}

// must calls log.Fatal if the error is non-nil.
func must(err error, msg string) {
	if err != nil {
		log.Fatalln(msg + ": " + err.Error())
	}
}
