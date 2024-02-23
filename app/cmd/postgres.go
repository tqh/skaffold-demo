package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

const (
	createTable = `
	CREATE TABLE IF NOT EXISTS 
		counters (
		name text PRIMARY KEY,
		value INTEGER 
	)
`
	incrementCounter = `
	INSERT INTO counters(name, value)
	VALUES ('%s', %d)
	ON CONFLICT(name) DO UPDATE 
		SET value=counters.value + EXCLUDED.value
	RETURNING value;
`
)

var dbpool *pgxpool.Pool

func mustExec[T any](ctx context.Context, fn func(context.Context) (T, error)) T {
	for range 20 {
		if result, err := fn(ctx); err == nil {
			return result
		}
		time.Sleep(4 * time.Second)
	}
	panic("failed to init postgres")
}

func init() {
	ctx := context.Background()
	dburl := os.Getenv("DATABASE_URL")
	dbpool = mustExec(ctx, func(ctx context.Context) (*pgxpool.Pool, error) { return pgxpool.New(ctx, dburl) })
	info := mustExec(ctx, func(ctx context.Context) (pgconn.CommandTag, error) { return dbpool.Exec(ctx, createTable) })
	log.Printf("table creation affected %d rows\n", info.RowsAffected())
}

func pgcount(name string) int64 {
	var count int64
	row := dbpool.QueryRow(context.Background(), fmt.Sprintf(incrementCounter, name, 1))
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err, "failed to increment counter")
	}
	return count
}
