// Package postgres provides a Postgres dependency for selftest, either from the
// host machine or from an ephemeral docker container.
package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func getPool(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(URL())
	if err != nil {
		return nil, err
	}
	return pgxpool.ConnectConfig(ctx, cfg)
}

func URL() string {
	return os.Getenv("TODO_API_TEST_DATABASE_URL")
}

func GetPool() (*pgxpool.Pool, error) {
	return getPool(context.Background())
}

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := pgx.Connect(ctx, URL())
	return err
}

// Reset truncates all tables, except for migrations.
func Reset() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := getPool(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	rows, _ := conn.Query(ctx, `
		select table_name
		from information_schema.tables
		where table_schema = 'public'
		  and table_name != 'migrations';
	`)
	cmd := ""
	for rows.Next() {
		var t string
		_ = rows.Scan(&t)
		cmd += fmt.Sprintf("truncate table %s;\n", t)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	_, err = conn.Exec(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func Truncate(table string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := getPool(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(ctx, fmt.Sprintf("truncate table %s;", table))
	return err
}
