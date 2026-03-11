package config

import (
	"embed"
	"fmt"
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
	SessionName   string // `mapstructure:"session_name"`
	SessionHours  int    // `mapstructure:"session_hours"`

	AuthEnabled string // `mapstructure:"auth_google_id"`

	AuthGoogleID     string // `mapstructure:"auth_google_id"`
	AuthGoogleSecret string // `mapstructure:"auth_google_secret"`

	AuthGitlabID     string // `mapstructure:"auth_gitlab_id"`
	AuthGitlabSecret string // `mapstructure:"auth_gitlab_secret"`

	AuthGiteaID     string // `mapstructure:"auth_gitea_id"`
	AuthGiteaSecret string // `mapstructure:"auth_gitea_secret"`

	TlsPrivateKey string // `mapstructure:"tls_privkey"`
	TlsChain      string // `mapstructure:"tls_chain"`

	Menus []Menu
}

func New() (*Config, error) {
	cfg := &Config{}

	// Load from environment with defaults
	pageSize := 10
	if ps := os.Getenv("PAGE_SIZE"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}
	cfg.PageSize = pageSize

	pageCount := 0
	if pc := os.Getenv("PAGE_COUNT"); pc != "" {
		fmt.Sscanf(pc, "%d", &pageCount)
	}
	cfg.PageCount = pageCount

	sessionName := "session"
	if sn := os.Getenv("SESSION_NAME"); sn != "" {
		sessionName = sn
	}
	cfg.SessionName = sessionName

	sessionHours := 24
	if sh := os.Getenv("SESSION_HOURS"); sh != "" {
		fmt.Sscanf(sh, "%d", &sessionHours)
	}
	cfg.SessionHours = sessionHours

	cfg.DbUri = os.Getenv("DB_URI")
	cfg.WebPrefix = os.Getenv("WEB_PREFIX")
	cfg.WebListener = os.Getenv("WEB_LISTENER")
	cfg.AuthEnabled = os.Getenv("AUTH_ENABLED")
	cfg.AuthGoogleID = os.Getenv("AUTH_GOOGLE_ID")
	cfg.AuthGoogleSecret = os.Getenv("AUTH_GOOGLE_SECRET")
	cfg.AuthGitlabID = os.Getenv("AUTH_GITLAB_ID")
	cfg.AuthGitlabSecret = os.Getenv("AUTH_GITLAB_SECRET")
	cfg.AuthGiteaID = os.Getenv("AUTH_GITEA_ID")
	cfg.AuthGiteaSecret = os.Getenv("AUTH_GITEA_SECRET")
	cfg.TlsChain = os.Getenv("TLS_CHAIN")
	cfg.TlsPrivateKey = os.Getenv("TLS_PRIVKEY")

	if err := checkStringNotEmpty(cfg.DbUri, "DB_URI"); err != nil {
		return nil, err
	}
	if err := checkStringNotEmpty(cfg.WebListener, "WEB_LISTENER"); err != nil {
		return nil, err
	}
	// checkStringNotEmpty(cfg.AuthGoogleID, "AUTH_GOOGLE_ID")
	// checkStringNotEmpty(cfg.AuthGoogleSecret, "AuthGoogleSecret")
	// checkStringNotEmpty(cfg.AuthGitlabID, "AUTH_GITLAB_ID")
	// checkStringNotEmpty(cfg.AuthGitlabSecret, "AuthGitlabSecret")
	return cfg, nil
}

func checkStringNotEmpty(s string, text string) error {
	if len(s) < 1 {
		return fmt.Errorf("required config not set: %s", text)
	}
	return nil
}

// --- EOF
