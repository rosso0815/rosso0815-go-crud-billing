package services

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rosso0815/rosso0815-go-crud-billing/config"
	"github.com/rosso0815/rosso0815-go-crud-billing/db"
)

type Store struct {
	Db  *db.Db
	Cfg *config.Config
}

// func NewStore(db *pgxpool.Pool, queries *db_gen.Queries, cfg *config.Config) *PgxStore {
func NewStore(pool *pgxpool.Pool, cfg *config.Config) *Store {
	pgxStore := Store{}
	pgxStore.Db = db.NewDb(pool, cfg)
	pgxStore.Cfg = cfg
	return &pgxStore
}

func (m *Store) Close() {
	m.Db.Close()
}

// --- EOF
