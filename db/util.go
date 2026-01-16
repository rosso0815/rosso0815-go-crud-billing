package db

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func LoadSQLFile(dbpool *pgxpool.Pool, sqlFile string) error {
	file, err := os.ReadFile(sqlFile)
	if err != nil {
		return errors.New(("cannot open file"))
	}
	_, err = dbpool.Exec(context.Background(), string(file))
	return err
}

// ---EOF
