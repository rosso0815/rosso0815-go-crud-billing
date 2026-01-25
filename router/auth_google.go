// Package router tbd
package router

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	config "github.com/rosso0815/rosso0815-go-crud-billing/config"
)

var googleOauthConfig *oauth2.Config

const (
	oauthGoogleURLAPI   = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	google_url_login    = "/auth/google/login"
	google_url_callback = "/auth/google/callback"
)

// ProfileDataGoogle tbd
type ProfileDataGoogle struct {
	Email string
}

func setupOauthGoogle(cfg *config.Config) {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://%s%s%s", cfg.WebListener, cfg.WebPrefix, google_url_callback),
		ClientID:     cfg.AuthGoogleID,
		ClientSecret: cfg.AuthGoogleSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (rt *Router) oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w, r, rt.sessionManager)
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (rt *Router) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	oauthState := rt.sessionManager.Get(r.Context(), "oauthstate")
	if r.FormValue("state") != oauthState {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data, err := getUserDataFromGoogle(r.FormValue("code"))
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

func generateStateOauthCookie(w http.ResponseWriter, r *http.Request, sessionManager *scs.SessionManager) string {
	log.Println("generateStateOauthCookie")
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err.Error())
	}
	state := base64.URLEncoding.EncodeToString(b)
	log.Println("state", state)
	sessionManager.Put(r.Context(), "oauthstate", state)
	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleURLAPI + token.AccessToken)
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
