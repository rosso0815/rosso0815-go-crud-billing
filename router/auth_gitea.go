// Package router tbd
package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
)

var (
	giteaOauthConfig *oauth2.Config
)

const (
	oauthGiteaURLAPI   = "https://gitea.com/user/settings/applications"
	gitea_url_login    = "/auth/gitea/login"
	gitea_url_callback = "/auth/gitea/callback"
)

func setupOauthGitea(cfg *config.Config) {
	giteaOauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://%s%s%s", cfg.WebListener, cfg.WebPrefix, gitea_url_callback),
		ClientID:     cfg.AuthGiteaID,
		ClientSecret: cfg.AuthGoogleSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     oauth2.Endpoint{},
	}
}

func (rt *Router) oauthGiteaLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w, r, rt.sessionManager)
	u := giteaOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (rt *Router) oauthGiteaCallback(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w, r, rt.sessionManager)
	if r.FormValue("state") != oauthState {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data, err := getUserDataFromGitea(r.FormValue("code"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var profileData ProfileDataGoogle
	err = json.Unmarshal(data, &profileData)
	if err != nil {
		log.Println(err.Error())
	}
	rt.sessionManager.Put(r.Context(), "userid", profileData.Email)
	redirect_url := rt.sessionManager.GetString(r.Context(), "redirect_url")
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	fullUrlPath := fmt.Sprintf("%v://%v%v", scheme, r.Host, redirect_url)
	http.Redirect(w, r, fullUrlPath, http.StatusTemporaryRedirect)
}

func getUserDataFromGitea(code string) ([]byte, error) {
	token, err := giteaOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGiteaURLAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

/// ---- EOF
