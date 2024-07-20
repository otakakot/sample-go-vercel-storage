package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	pool, err := GetPool(ctx, GetDSN())
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	row, err := pool.Query(ctx, "SELECT 1")
	if err != nil {
		panic(err)
	}

	defer row.Close()

	for row.Next() {
		var one int
		if err := row.Scan(&one); err != nil {
			panic(err)
		}

		println(one)
	}
}

func GetDSN() string {
	return "postgres://default:" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOST") + ":5432/" + os.Getenv("POSTGRES_DATABASE") + "?sslmode=require"
}

func GetPool(
	ctx context.Context,
	dsn string,
) (*pgxpool.Pool, error) {
	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, conn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
