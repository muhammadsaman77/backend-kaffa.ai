package configs

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func InitDatabase(ctx context.Context) *pgxpool.Pool {
	dbUrl := viper.GetString("db_url")
	fmt.Println("Connecting to database at:", dbUrl)
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database!")

	return pool
}
