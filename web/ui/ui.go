package ui

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"

	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
	router "github.com/rosso0815/rosso0815-go-crud-billing/router"
)

type Message int

const (
	Debug Message = iota + 1
	Normal
	Alert
)

type Action int

const (
	Add Action = iota + 1
	Edit
	Delete
	List
)

// Link tbd
type Link struct {
	Href templ.SafeURL
	Text string
}

// LayoutData tbd
type LayoutData struct {
	Status  string
	Prefix  string
	Title   string
	Actions []templ.Component
	Search  templ.Component
	// Debug   string
}

// CrudList tbd
type CrudList struct {
	Base
	Header []templ.Component
	Items  []templ.Component
	Count  int
}

type PaginationLink struct {
	URL  templ.SafeURL
	Text string
}

type Ui struct {
	Base
	cfg *config.Config
}

func New(sessionManager *scs.SessionManager, cfg *config.Config) *Ui {
	u := Ui{}
	u.SessionManager = sessionManager
	u.cfg = cfg
	router.RegisterRoute("GET "+u.cfg.WebPrefix+"/auth/login", u.AuthLogin)
	return &u
}

// AuthLogin tbd
func (u *Ui) AuthLogin(w http.ResponseWriter, r *http.Request) {
	links := []Link{
		{Text: "Gitlab", Href: templ.URL(u.cfg.WebPrefix + "/auth/gitlab/login")},
		{Text: "Google", Href: templ.URL(u.cfg.WebPrefix + "/auth/google/login")},
		{Text: "Gitea", Href: templ.URL(u.cfg.WebPrefix + "/auth/gitea/login")},
	}
	u.Update(w, r)
	err := authLogin(&u.Base, links, u.cfg).Render(r.Context(), w)
	if err != nil {
		log.Fatalln("ui", "authLogin", err)
	}
}

// --- EOF
