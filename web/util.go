package web

import (
	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
	services "github.com/rosso0815/rosso0815-go-crud-billing/services"
	ui "github.com/rosso0815/rosso0815-go-crud-billing/web/ui"
)

type Web struct {
	ui.Base
	store *services.Store
	cfg   *config.Config
}
