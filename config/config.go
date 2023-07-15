package config

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v9"
	"golang.org/x/xerrors"
)

type Portal struct {
	// APIルートURL
	API   string `env:"PORTAL_API,required"`
	Token string `env:"PORTAL_TOKEN,required"`
}

type Auth0 struct {
	Domain       string `env:"AUTH0_DOMAIN,required"`
	ClientID     string `env:"AUTH0_CLIENT_ID,required"`
	ClientSecret string `env:"AUTH0_CLIENT_SECRET,required"`
}

type Config struct {
	Portal Portal
	Auth0  Auth0
}

func ReadConfig() (*Config, error) {
	var cfg Config

	err := env.Parse(&cfg.Portal)
	if err != nil {
		return nil, xerrors.Errorf("failed to perse config: %w", err)
	}

	err = env.Parse(&cfg.Auth0)
	if err != nil {
		return nil, xerrors.Errorf("failed to perse config: %w", err)
	}

	fmt.Println(cfg)

	portalAPIURL, err := url.Parse(cfg.Portal.API)

	if err != nil {
		return nil, fmt.Errorf("envvar \"PORTAL_API\" must be a valid URL that indicates API Root of MIS.W Portal: %w", err)
	}

	if !strings.HasSuffix(portalAPIURL.Path, "/") {
		log.Println("there is no \"/\" at the end of \"PORTAL_API\". you may have problems around relative paths calculation.")
	}

	return &cfg, err
}
