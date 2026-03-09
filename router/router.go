// Package router tbd
package router

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
)

type Router struct {
	mux            *http.ServeMux
	sessionManager *scs.SessionManager
	staticFs       embed.FS
	cfg            *config.Config
	routes         []route
}

type route struct {
	path    string
	handler http.HandlerFunc
}

// Router type now holds routes instead of using global state
func New(staticFs embed.FS, cfg *config.Config, pool *pgxpool.Pool) *Router {
	r := Router{}
	r.cfg = cfg
	r.staticFs = staticFs
	r.mux = http.NewServeMux()
	r.routes = make([]route, 0) // Initialize routes slice
	setupOauthGoogle(r.cfg)
	setupOauthGitlab(r.cfg)
	setupOauthGitea(r.cfg)
	r.sessionManager = scs.New()
	r.sessionManager.Store = pgxstore.New(pool)
	r.sessionManager.Lifetime = time.Duration(r.cfg.SessionHours) * time.Hour
	r.sessionManager.Cookie.Name = r.cfg.SessionName
	staticFsJs, _ := fs.Sub(staticFs, "static/js")
	r.GetMux().Handle(cfg.WebPrefix+"/js/", http.StripPrefix(cfg.WebPrefix+"/js", http.FileServerFS(staticFsJs)))
	staticFsCss, _ := fs.Sub(staticFs, "static/css")
	r.GetMux().Handle(cfg.WebPrefix+"/css/", http.StripPrefix(cfg.WebPrefix+"/css", http.FileServerFS(staticFsCss)))
	staticFsFonts, _ := fs.Sub(staticFs, "static/fonts")
	r.GetMux().Handle(cfg.WebPrefix+"/fonts/", http.StripPrefix(cfg.WebPrefix+"/fonts", http.FileServerFS(staticFsFonts)))
	return &r
}

// RegisterRoute adds a route to the router
func (rt *Router) RegisterRoute(path string, handler http.HandlerFunc) {
	rt.routes = append(rt.routes, route{
		path:    path,
		handler: handler,
	})
}

func (rt *Router) SetupRoutes() {
	rt.RegisterRoute("GET "+rt.cfg.WebPrefix+google_url_login, rt.oauthGoogleLogin)
	rt.RegisterRoute("GET "+rt.cfg.WebPrefix+google_url_callback, rt.oauthGoogleCallback)
	rt.RegisterRoute("GET "+rt.cfg.WebPrefix+gitlab_url_login, rt.oauthGitlabLogin)
	rt.RegisterRoute("GET "+rt.cfg.WebPrefix+gitlab_url_callback, rt.oauthGitlabCallback)
	rt.RegisterRoute("GET "+rt.cfg.WebPrefix+gitea_url_login, rt.oauthGiteaLogin)
	rt.RegisterRoute("GET "+rt.cfg.WebPrefix+gitea_url_callback, rt.oauthGiteaCallback)
	rt.RegisterRoute("GET "+rt.cfg.WebPrefix+"/status", rt.StatusHandler)
	rt.RegisterRoute("GET "+rt.cfg.WebPrefix+"", rt.redirectToCustomer)
	for _, r := range rt.routes {
		if rt.cfg.AuthEnabled == "true" {
			rt.mux.Handle(r.path, rt.sessionManager.LoadAndSave(rt.AuthLoginHandler(logger(r.handler), rt.cfg)))
		} else {
			rt.mux.Handle(r.path, rt.sessionManager.LoadAndSave(logger(r.handler)))
		}
	}
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("logger->request :", r.Method, r.URL.Path, r.URL.Query())
			next.ServeHTTP(w, r)
		})
}

func (rt *Router) StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "r.tls: %v \n", r.TLS)
	fmt.Fprintf(w, "r.URL.Scheme: %s \n", r.URL.Scheme)
	fmt.Fprintf(w, "r.Host: %s \n", r.Host)
	fmt.Fprintf(w, "r.URL.Path: %s \n", r.URL.Path)
}

func (rt *Router) redirectToCustomer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, rt.cfg.WebPrefix+"/ui/customer", http.StatusSeeOther)
}

func (r *Router) GetSessionManager() *scs.SessionManager {
	return r.sessionManager
}

func (r *Router) GetMux() *http.ServeMux {
	return r.mux
}

// --- EOF
