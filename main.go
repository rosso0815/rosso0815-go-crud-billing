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

//go:embed config.yaml
var embedConfig embed.FS

//go:embed db/schema/*.sql
var embedMigrations embed.FS

func cleanup() {
	os.Exit(1)
}

func web_run(cfg *config.Config) {
	ctx := context.Background()

	store := services.NewStore(ctx, cfg)
	defer store.Close()

	r := router.New(embedStatic, cfg)
	ui.New(r.GetSessionManager(), cfg)
	web.NewCustomer(store, r.GetSessionManager(), cfg)
	web.NewInvoice(store, r.GetSessionManager(), cfg)
	// web.NewUserkv(store, r.GetSessionManager(), cfg)
	r.SetupRoutes()
	fmt.Printf("running on http://%s%s\n", cfg.WebListener, cfg.WebPrefix)
	srv := &http.Server{
		Handler:      r.GetMux(),
		Addr:         cfg.WebListener,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
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

	cfg := config.New(&embedConfig)

	var action string = "web"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}
	switch action {
	case "web":
		// services.NewScheduler()
		web_run(cfg)
	case "db_migrate_down":
		db.Db_migration_up(&embedMigrations, cfg)
	case "db_migrate_up":
		db.Db_migration_down(&embedMigrations, cfg)
	}
}

// --- EOF
