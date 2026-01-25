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
	"golang.org/x/oauth2/gitlab"

	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
)

var gitlabOauthConfig *oauth2.Config

const (
	gitlab_url_login    = "/auth/gitlab/login"
	gitlab_url_callback = "/auth/gitlab/callback"
	oauthGitlabUrlAPI   = "https://gitlab.com/api/v4/user?access_token="
)

type ProfileDataGitlab struct {
	Email string
}

func setupOauthGitlab(cfg *config.Config) {
	gitlabOauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://%s%s%s", cfg.WebListener, cfg.WebPrefix, gitlab_url_callback),
		ClientID:     cfg.AuthGitlabID,
		ClientSecret: cfg.AuthGitlabSecret,
		Scopes:       []string{"read_user", "openid", "profile", "email"},
		Endpoint:     gitlab.Endpoint,
	}
}

func (rt *Router) oauthGitlabLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w, r, rt.sessionManager)
	u := gitlabOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (rt *Router) oauthGitlabCallback(w http.ResponseWriter, r *http.Request) {
	oauthState := rt.sessionManager.Get(r.Context(), "oauthstate")
	if r.FormValue("state") != oauthState {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data, err := getUserDataFromGitlab(r.FormValue("code"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var profileData ProfileDataGitlab
	err = json.Unmarshal(data, &profileData)
	if err != nil {
		log.Println("auth_gitlab", "error", err.Error())
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

func getUserDataFromGitlab(code string) ([]byte, error) {
	token, err := gitlabOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGitlabUrlAPI + token.AccessToken)
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
