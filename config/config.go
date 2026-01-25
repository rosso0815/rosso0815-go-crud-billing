package config

import (
	"embed"
	"fmt"
	"log"
	"os"
)

type Menu struct {
	Name string
	Path string
}

type Config struct {
	embedFs *embed.FS

	// can be configured by env GOAPP_DB_URI=
	DbUri         string // `mapstructure:"db_uri"`
	WebListener   string // `mapstructure:"web_listener"`
	WebPrefix     string // `mapstructure:"web_prefix"`
	PageSize      int    // `mapstructure:"page_size"`
	PageCount     int    // `mapstructure:"page_count"`
	SessionSecret string // `mapstructure:"session_secret"`

	AuthGoogleID     string // `mapstructure:"auth_google_id"`
	AuthGoogleSecret string // `mapstructure:"auth_google_secret"`

	AuthGitlabID     string // `mapstructure:"auth_gitlab_id"`
	AuthGitlabSecret string // `mapstructure:"auth_gitlab_secret"`

	AuthGiteaID     string // `mapstructure:"auth_gitea_id"`
	AuthGiteaSecret string // `mapstructure:"auth_gitea_secret"`

	Menus []Menu
}

func New(e *embed.FS) *Config {
	cfg := &Config{}
	cfg.embedFs = e
	cfg.PageSize = 10
	cfg.PageCount = 0

	cfg.DbUri = os.Getenv("DB_URI")
	fmt.Println("DB_URI: ", cfg.DbUri)

	cfg.WebPrefix = os.Getenv("WEB_PREFIX")
	fmt.Println("WEB_PREFIX: ", cfg.WebPrefix)

	cfg.WebListener = os.Getenv("WEB_LISTENER")
	fmt.Println("WEB_LISTENER: ", cfg.WebListener)

	cfg.AuthGoogleID = os.Getenv("AUTH_GOOGLE_ID")
	fmt.Println("AUTH_GOOGLE_ID: ", cfg.AuthGoogleID)

	cfg.AuthGoogleSecret = os.Getenv("AUTH_GOOGLE_SECRET")
	fmt.Println("AUTH_GOOGLE_SECRET: ", cfg.AuthGoogleSecret)

	cfg.AuthGitlabID = os.Getenv("AUTH_GITLAB_ID")
	fmt.Println("AUTH_GITLAB_ID: ", cfg.AuthGitlabID)

	cfg.AuthGitlabSecret = os.Getenv("AUTH_GITLAB_SECRET")
	fmt.Println("AUTH_GITLAB_SECRET: ", cfg.AuthGitlabSecret)

	cfg.AuthGiteaID = os.Getenv("AUTH_GITEA_ID")
	fmt.Println("AUTH_GITEA_ID: ", cfg.AuthGiteaID)

	cfg.AuthGiteaSecret = os.Getenv("AUTH_GITEA_SECRET")
	fmt.Println("AUTH_GITEA_SECRET: ", cfg.AuthGiteaSecret)

	checkStringNotEmpty(cfg.DbUri, "DB_URI")
	checkStringNotEmpty(cfg.WebListener, "WEB_LISTENER")
	checkStringNotEmpty(cfg.AuthGoogleID, "AUTH_GOOGLE_ID")
	checkStringNotEmpty(cfg.AuthGoogleSecret, "AuthGoogleSecret")
	checkStringNotEmpty(cfg.AuthGitlabID, "AUTH_GITLAB_ID")
	checkStringNotEmpty(cfg.AuthGitlabSecret, "AuthGitlabSecret")
	return cfg
}

func checkStringNotEmpty(s string, text string) {
	if len(s) < 1 {
		log.Panicln(text)
	}
}

// --- EOF
