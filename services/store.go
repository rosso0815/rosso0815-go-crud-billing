package services

import (
	"context"

	"github.com/rosso0815/rosso0815-go-crud-billing/config"
	"github.com/rosso0815/rosso0815-go-crud-billing/db"
)

type Store struct {
	Db  *db.Db
	Cfg *config.Config
}

// func NewStore(db *pgxpool.Pool, queries *db_gen.Queries, cfg *config.Config) *PgxStore {
func NewStore(ctx context.Context, cfg *config.Config) *Store {
	pgxStore := Store{}
	pgxStore.Db = db.NewDb(ctx, cfg)
	pgxStore.Cfg = cfg
	return &pgxStore
}

func (m *Store) Close() {
	m.Db.Close()
}

// --- EOF
