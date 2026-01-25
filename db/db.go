package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rosso0815/rosso0815-go-crud-billing/config"
	db_gen "github.com/rosso0815/rosso0815-go-crud-billing/db/generated"
)

type Db struct {
	Db      *pgxpool.Pool
	Queries *db_gen.Queries
	Cfg     *config.Config
}

func NewDb(ctx context.Context, cfg *config.Config) *Db {
	pool, err := pgxpool.New(ctx, cfg.DbUri)
	if err != nil {
		log.Fatalln("cannot connect", err)
	}
	pgxStore := Db{}
	pgxStore.Db = pool
	pgxStore.Cfg = cfg
	pgxStore.Queries = db_gen.New()
	return &pgxStore
}

func (db *Db) Close() {
	db.Db.Close()
}

// --- EOF
