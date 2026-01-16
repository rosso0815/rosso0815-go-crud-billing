// Package userkv provides simple per-user key/value management UI and handlers.
package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"

	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
	router "github.com/rosso0815/rosso0815-go-crud-billing/router"
	store "github.com/rosso0815/rosso0815-go-crud-billing/services"
	ui "github.com/rosso0815/rosso0815-go-crud-billing/web/ui"
)

func NewUserkv(store *store.Store, sessionManager *scs.SessionManager, cfg *config.Config) *Web {
	b := Web{}
	b.store = store
	b.Path = "userkv"
	b.SessionManager = sessionManager
	b.cfg = cfg
	b.AddAction = ui.CrudActionAdd("tbd")
	cfg.Menus = append(cfg.Menus, config.Menu{Name: "Key Values", Path: fmt.Sprintf("ui/%s", b.Path)})
	router.RegisterRoute(fmt.Sprintf("GET %s/ui/%s", b.cfg.WebPrefix, b.Path), b.ListAll)
	return &b
}

// ListAll shows all keys for the current session user.
func (m *Web) ListAll(w http.ResponseWriter, r *http.Request) {
	m.Update(w, r)
	var components []templ.Component
	var header []templ.Component = []templ.Component{
		ui.CrudHeaderSort("UserName", m.cfg),
		ui.CrudHeaderSort("Key", m.cfg),
		ui.CrudHeaderSort("Value", m.cfg)}
	if m.User == "" {
		m.Base.MessageType = ui.Alert
		m.Base.MessageText = "not logged in"
		ui.CrudMessageOnly(m.Base, m.cfg).Render(r.Context(), w)
		return
	}
	items, err := m.store.UserkvList(r.Context(), m.User)
	if err != nil {
		log.Println("err", err.Error())
		return
	}
	count := len(items)
	for _, item := range items {
		l := userkvCrudLine(item, m.cfg)
		components = append(components, l)
	}
	ui.CrudTable(
		ui.CrudList{
			Base:   m.Base,
			Header: header,
			Items:  components,
			Count:  count,
		}, m.cfg).Render(r.Context(), w)

}
