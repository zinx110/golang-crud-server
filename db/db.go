package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DBpool *pgxpool.Pool

func InitDb() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Unable to parse database URL: %v", err)
	}
	config.MaxConns = 10
	config.MaxConnLifetime = time.Hour

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DBpool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Unable to create conntection pool: %v", err)
	}

	if err := DBpool.Ping(ctx); err != nil {
		log.Fatal("Unable to connnect to database: %v", err)
	}

	log.Println("Database conection pool created successfully")

}
