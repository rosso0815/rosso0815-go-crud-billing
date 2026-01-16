package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func Test_ShowInvoices(t *testing.T) {
	t.Log("--> ShowInvoices")

	dbConfig, err := pgx.ParseConfig(os.Getenv("GOAPP_DB_URI"))
	if err != nil {
		log.Fatal("Unable to parse connString: ", err)
	}

	dbCon, err := pgx.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	defer dbCon.Close(context.Background())

	result, err := dbCon.Query(context.Background(), "select now()")
	if err != nil {
		log.Fatal("Unable to execute query: ", err)
	}
	log.Println("Database result", result)
	t.Log("<-- ShowInvoices")
}
