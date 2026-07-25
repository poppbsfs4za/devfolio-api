package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/example/devfolio-api/internal/config"
	"github.com/example/devfolio-api/internal/infrastructure/persistence/gormmodel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// buildDSN prefers an explicit DATABASE_URL (e.g. Neon's
// "postgres://user:pass@host/db?sslmode=require" connection string) and
// falls back to the discrete DB_* fields for local/Docker Compose use.
// It never logs the resulting DSN since it may contain credentials.
func buildDSN(cfg config.DBConfig) string {
	if cfg.DatabaseURL != "" {
		return cfg.DatabaseURL
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
		cfg.TimeZone,
	)
}

// NewPostgres opens a GORM/Postgres connection and applies a conservative
// connection pool suited to Neon's free-tier compute/connection limits.
func NewPostgres(cfg config.DBConfig, appEnv string) (*gorm.DB, error) {
	dsn := buildDSN(cfg)

	// Keep GORM's own SQL logging quiet outside local dev - it is noisy and
	// not needed for a low-traffic personal site; request-level logging is
	// handled by the Fiber logger middleware instead.
	logLevel := logger.Warn
	if appEnv == "local" || appEnv == "" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		// Do not include the DSN/err verbatim if it might echo credentials;
		// Postgres/pgx connection errors normally don't include the password,
		// but we avoid printing cfg fields here regardless.
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB: %w", err)
	}

	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 4
	}
	maxIdle := cfg.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 2
	}
	connMaxLifetime := cfg.ConnMaxLifetime
	if connMaxLifetime <= 0 {
		connMaxLifetime = 30
	}
	connMaxIdleTime := cfg.ConnMaxIdleTime
	if connMaxIdleTime <= 0 {
		connMaxIdleTime = 5
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(connMaxIdleTime) * time.Minute)

	log.Printf("[database] connected (pool: max_open=%d max_idle=%d conn_max_lifetime=%dm conn_max_idle_time=%dm)",
		maxOpen, maxIdle, connMaxLifetime, connMaxIdleTime)

	return db, nil
}

// ConnectWithRetry attempts to connect a bounded number of times with a fixed
// backoff. This is meant to smooth over brief, transient connection blips at
// cold start (e.g. Neon's compute waking from auto-suspend), not to wait out
// a hard/prolonged outage such as an exhausted monthly compute quota - if all
// attempts fail, the caller is expected to continue starting the server in a
// degraded state rather than crash the process.
func ConnectWithRetry(cfg config.DBConfig, appEnv string, attempts int, delay time.Duration) (*gorm.DB, error) {
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for i := 1; i <= attempts; i++ {
		db, err := NewPostgres(cfg, appEnv)
		if err == nil {
			return db, nil
		}
		lastErr = err
		log.Printf("[database] connection attempt %d/%d failed: %v", i, attempts, err)
		if i < attempts {
			time.Sleep(delay)
		}
	}
	return nil, lastErr
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("[database] running automigrate...")
	return db.AutoMigrate(
		&gormmodel.User{},
		&gormmodel.Profile{},
		&gormmodel.Project{},
		&gormmodel.Tag{},
		&gormmodel.Post{},
		&gormmodel.PostTag{},
	)
}

// Status tracks whether a database connection is currently available and is
// shared between the health/readiness handlers and the DB-guard middleware.
// A nil DB (connection failed at startup) is a valid, non-crashing state.
type Status struct {
	db *gorm.DB
}

func NewStatus(db *gorm.DB) *Status {
	return &Status{db: db}
}

// Available reports whether a DB handle exists at all (cheap, no I/O).
// Used by the request-time guard middleware to short-circuit DB-dependent
// routes with a 503 instead of letting them panic/error deeper in the stack.
func (s *Status) Available() bool {
	return s != nil && s.db != nil
}

// Ping performs a real round-trip to confirm the database is reachable, used
// by the /ready(z) endpoint. It returns a descriptive error rather than
// panicking so callers can produce a structured 503 response.
func (s *Status) Ping(ctx context.Context) error {
	if s == nil || s.db == nil {
		return errors.New("database not connected")
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("get underlying sql.DB: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}
