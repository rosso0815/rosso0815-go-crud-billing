package router

import (
	"net/http"
	"strings"

	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
)

func (rt *Router) AuthLoginHandler(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.RequestURI, cfg.WebPrefix+"/auth") ||
			strings.HasPrefix(req.RequestURI, cfg.WebPrefix+"/api") {
			next.ServeHTTP(w, req)
			return
		}
		user := rt.sessionManager.GetString(req.Context(), "userid")
		if len(user) < 1 {
			rt.sessionManager.Put(req.Context(), "redirect_url", req.URL.Path)
			http.Redirect(w, req, cfg.WebPrefix+"/auth/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, req)
	})
}

// --- EOF
