package ui

import (
	"fmt"
	// "log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
)

type BaseCrud struct {
	PageSize  int
	PageCount int
	Search    string
	SortName  string
	SortOrder string
}

type Base struct {
	BaseCrud
	SessionManager *scs.SessionManager
	AddAction      templ.Component
	Path           string // endpath of url
	MessageType    Message
	MessageText    string
	User           string
	Title          string
	HxRequest      bool
	FullUrlPath    string // http:.....
}

func (b *Base) Update(w http.ResponseWriter, r *http.Request) {
	b.PageSize = 0
	pagesize_s := b.SessionManager.GetInt(r.Context(), "pagesize")
	pagesize_q := r.URL.Query().Get("pagesize")
	if len(pagesize_q) < 1 && pagesize_s < 1 {
		b.SessionManager.Put(r.Context(), "pagesize", 5)
		b.PageSize = 5
	} else if len(pagesize_q) > 0 {
		b.PageSize, _ = strconv.Atoi(pagesize_q)
		b.SessionManager.Put(r.Context(), "pagesize", b.PageSize)
	} else if pagesize_s > 0 {
		b.PageSize = pagesize_s
	}

	b.PageCount = 0
	pagecount_s := r.URL.Query().Get("pagecount")
	if len(pagecount_s) > 0 {
		b.PageCount, _ = strconv.Atoi(pagecount_s)
	}

	b.User = b.SessionManager.GetString(r.Context(), "userid")

	b.Search = r.URL.Query().Get("search")

	if r.Header.Get("HX-Request") == "true" {
		b.HxRequest = true
	} else {
		b.HxRequest = false
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	b.FullUrlPath = fmt.Sprintf("%v://%v%v", scheme, r.Host, r.URL.String())
}

// --- EOF
