package db

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Open connects GORM to PostgreSQL. Every entry point (the API, cmd/*) opens the
// database through here so they all share the same pgx settings.
//
// pgx caches named prepared statements by default (QueryExecModeCacheStatement).
// Behind Neon's pooler (PgBouncer in transaction mode) that breaks after schema
// changes: a statement prepared before ALTER TABLE fails with "cached plan must not
// change result type" (0A000), and statement names clash between clients sharing a
// server connection ("prepared statement name is already in use", 08P01).
//
// QueryExecModeExec uses only the unnamed statement: parse, bind and execute go in one
// round trip, so nothing is cached between queries. Parameters are sent as text with
// types inferred by the server. Consequence for driver.Valuer types: a []byte value is
// always encoded as bytea, so a Valuer for a json/jsonb column must return a string.
func Open(dsn string, cfg *gorm.Config) (*gorm.DB, error) {
	pgxCfg, err := pgxConfig(dsn)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &gorm.Config{}
	}
	return gorm.Open(postgres.New(postgres.Config{Conn: stdlib.OpenDB(*pgxCfg)}), cfg)
}

func pgxConfig(dsn string) (*pgx.ConnConfig, error) {
	c, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err) // pgx redacts the password
	}
	c.DefaultQueryExecMode = pgx.QueryExecModeExec
	return c, nil
}
