// repositories provides functions for interacting with a PostgreSQL database.
package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romankravchuk/muerta/internal/pkg/config"
)

// Client is an interface for defining methods to interact with a PostgreSQL client.
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func connectWithTries(
	fn func() (*pgxpool.Pool, error),
	attemtps int,
	delay time.Duration,
) (*pgxpool.Pool, error) {
	var (
		err  error
		pool *pgxpool.Pool
	)
	for attemtps > 0 {
		pool, err = fn()
		if err != nil {
			time.Sleep(delay)
			attemtps--
			continue
		}
		return pool, nil
	}
	return nil, err
}

// New returns a new PostgreSQL client connection pool based on the provided config.
//
// The function attempts to connect to the database maxAttempts times with a delay of 5 seconds between attempts.
func New(
	ctx context.Context,
	maxAttempts int,
	cfg *config.Config,
) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	pool, err := connectWithTries(func() (*pgxpool.Pool, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return nil, err
		}
		return pool, nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		log.Fatalf("error do with tries postgresql: %v", err)
	}
	return pool, nil
}
