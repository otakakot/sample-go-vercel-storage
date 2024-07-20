package main

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		panic("POSTGRES_URL is required")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			println(err)
		}
	}()

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	rows, err := db.QueryContext(ctx, "SELECT 1")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var one int
		if err := rows.Scan(&one); err != nil {
			panic(err)
		}

		println(one)
	}
}
