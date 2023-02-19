package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type pgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// NewPgConfig создает новый инстанс конфига
func NewPgConfig(username string, password string, host string, port string, database string) *pgConfig {
	return &pgConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

func NewClient(ctx context.Context, maxAttempts int, maxDelay time.Duration, cfg *pgConfig) (db *sql.DB, err error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port, cfg.Database,
	)
	println(dsn)

	err = DoWithAttempts(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var openErr error
		db, openErr = sql.Open("postgres", dsn)
		if openErr != nil {
			return openErr
		}

		if pingErr := db.PingContext(ctx); pingErr != nil {
			log.Println("Failed to connect to postgres... Going to do the next attempt")
			db.Close()

			return pingErr
		}

		return nil
	}, maxAttempts, maxDelay)

	if err != nil {
		log.Fatal("All attempts are exceeded. Unable to connect to postgres")
	}

	return db, nil
}

// DoWithAttempts задает лимит на количество попыток подключения к базе
func DoWithAttempts(fn func() error, maxAttempts int, delay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--

			continue
		}

		return nil
	}

	return err
}
