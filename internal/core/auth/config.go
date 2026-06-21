package core_auth

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JWTSecret  string        `envconfig:"JWT_SECRET" required:"true"`
	TokenTTL   time.Duration `envconfig:"TOKEN_TTL" default:"168h"`
	CookieName string        `envconfig:"COOKIE_NAME" default:"todoapp_token"`
	// CookieSecure should be true in production (HTTPS-only cookie). It is
	// kept configurable so local development over plain HTTP keeps working.
	CookieSecure bool `envconfig:"COOKIE_SECURE" default:"false"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("AUTH", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get Auth config: %w", err))
	}

	return config
}
