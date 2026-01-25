// Package router tbd
package router

import (
	"context"
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
}

type route struct {
	path    string
	handler http.HandlerFunc
}

var routes []route

func RegisterRoute(path string, handler http.HandlerFunc) {
	routes = append(routes, route{
		path:    path,
		handler: handler,
	})
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("logger->request :", r.Method, r.URL.Path, r.URL.Query())
			next.ServeHTTP(w, r)
		})
}

func New(staticFs embed.FS, cfg *config.Config) *Router {
	r := Router{}
	r.cfg = cfg
	r.staticFs = staticFs
	r.mux = http.NewServeMux()
	setupOauthGoogle(r.cfg)
	setupOauthGitlab(r.cfg)
	setupOauthGitea(r.cfg)
	pool, err := pgxpool.New(context.Background(), r.cfg.DbUri)
	if err != nil {
		log.Fatal("router err", err)
	}
	r.sessionManager = scs.New()
	r.sessionManager.Store = pgxstore.New(pool)
	r.sessionManager.Lifetime = 24 * time.Hour
	r.sessionManager.Cookie.Name = "pfistera"
	staticFsJs, _ := fs.Sub(staticFs, "static/js")
	r.GetMux().Handle(cfg.WebPrefix+"/js/", http.StripPrefix(cfg.WebPrefix+"/js", http.FileServerFS(staticFsJs)))
	staticFsCss, _ := fs.Sub(staticFs, "static/css")
	r.GetMux().Handle(cfg.WebPrefix+"/css/", http.StripPrefix(cfg.WebPrefix+"/css", http.FileServerFS(staticFsCss)))
	staticFsFonts, _ := fs.Sub(staticFs, "static/fonts")
	r.GetMux().Handle(cfg.WebPrefix+"/fonts/", http.StripPrefix(cfg.WebPrefix+"/fons", http.FileServerFS(staticFsFonts)))
	return &r
}

func (rt *Router) SetupRoutes() {
	RegisterRoute("GET "+rt.cfg.WebPrefix+google_url_login, rt.oauthGoogleLogin)
	RegisterRoute("GET "+rt.cfg.WebPrefix+google_url_callback, rt.oauthGoogleCallback)
	RegisterRoute("GET "+rt.cfg.WebPrefix+gitlab_url_login, rt.oauthGitlabLogin)
	RegisterRoute("GET "+rt.cfg.WebPrefix+gitlab_url_callback, rt.oauthGitlabCallback)
	RegisterRoute("GET "+rt.cfg.WebPrefix+gitea_url_login, rt.oauthGiteaLogin)
	RegisterRoute("GET "+rt.cfg.WebPrefix+gitea_url_callback, rt.oauthGiteaCallback)
	RegisterRoute("GET "+rt.cfg.WebPrefix+"/status", rt.StatusHandler)
	RegisterRoute("GET "+rt.cfg.WebPrefix+"", rt.redirectToCustomer)
	for _, r := range routes {
		// rt.mux.Handle(r.path, rt.sessionManager.LoadAndSave(logger(r.handler)))
		rt.mux.Handle(r.path, rt.sessionManager.LoadAndSave(rt.AuthLoginHandler(logger(r.handler), rt.cfg)))
	}
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
