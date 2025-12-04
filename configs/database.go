package configs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitDatabase(ctx context.Context) *pgxpool.Pool {
	dbUrl := viper.GetString("db_url")
	fmt.Println("Connecting to database at:", dbUrl)
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		Log.Fatal("Unable to connect to database", zap.Error(err))

	}
	if err := pool.Ping(ctx); err != nil {
		Log.Fatal("Unable to ping database", zap.Error(err))

	}
	Log.Info("Connected to PostgreSQL database")
	return pool
}
