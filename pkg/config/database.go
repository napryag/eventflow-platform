package config

import (
	"database/sql"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/napryag/eventflow-platform/pkg/errs"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	PingTimeout     time.Duration
}

func (c DatabaseConfig) DSN() string {
	values := url.Values{}
	values.Set("sslmode", c.SSLMode)

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     "/" + c.Name,
		RawQuery: values.Encode(),
	}

	return u.String()
}

func (c DatabaseConfig) Validate() error {
	if strings.TrimSpace(c.Host) == "" {
		return errs.New("postgres host is required")
	}
	if strings.TrimSpace(c.Port) == "" {
		return errs.New("postgres port is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return errs.New("postgres database name is required")
	}
	if strings.TrimSpace(c.User) == "" {
		return errs.New("postgres user is required")
	}
	if strings.TrimSpace(c.Password) == "" {
		return errs.New("postgres password is required")
	}
	if strings.TrimSpace(c.SSLMode) == "" {
		return errs.New("postgres sslmode is required")
	}

	return nil
}

func ConfigurePool(db *sql.DB, cfg DatabaseConfig) {
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
}
