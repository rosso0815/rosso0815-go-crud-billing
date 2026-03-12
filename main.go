// Packahe main tbd
package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
	db "github.com/rosso0815/rosso0815-go-crud-billing/db"
	router "github.com/rosso0815/rosso0815-go-crud-billing/router"
	"github.com/rosso0815/rosso0815-go-crud-billing/services"
	web "github.com/rosso0815/rosso0815-go-crud-billing/web"
	ui "github.com/rosso0815/rosso0815-go-crud-billing/web/ui"
)

// FIXME: https://rollbar.com/blog/golang-error-logging-guide/
// DONE: slog back to log

//go:embed static/js/*js
//go:embed static/css/*css
//go:embed static/css/fonts/*woff
//go:embed static/css/fonts/*woff2
var embedStatic embed.FS

//go:embed db/schema/*.sql
var embedMigrations embed.FS

func cleanup() {
	os.Exit(1)
}

func web_run(cfg *config.Config) error {
	ctx := context.Background()

	// Create single database pool
	pool, err := pgxpool.New(ctx, cfg.DbUri)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	store := services.NewStore(pool, cfg)

	r := router.New(embedStatic, cfg, pool)
	ui.New(r.GetSessionManager(), cfg, r)
	web.NewCustomer(store, r.GetSessionManager(), cfg, r)
	web.NewInvoice(store, r.GetSessionManager(), cfg, r)
	r.SetupRoutes()
	log.Printf("running on http://%s%s\n", cfg.WebListener, cfg.WebPrefix)
	srv := &http.Server{
		Handler:      r.GetMux(),
		Addr:         cfg.WebListener,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// err = srv.ListenAndServe()
	fmt.Print("TLS_CHAIN", cfg.TlsChain)
	fmt.Print("TLS_PRIVKEY", cfg.TlsPrivateKey)
	err = srv.ListenAndServeTLS(cfg.TlsChain, cfg.TlsPrivateKey)
	if err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		cleanup()
		os.Exit(1)
	}()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var action string = "web"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}
	switch action {
	case "web":
		if err := web_run(cfg); err != nil {
			log.Fatalf("web_run failed: %v", err)
		}
	case "db_migrate_down":
		db.DbMigrationDown(&embedMigrations, cfg)
	case "db_migrate_up":
		db.DbMigrationUp(&embedMigrations, cfg)
	}
}

// --- EOF
