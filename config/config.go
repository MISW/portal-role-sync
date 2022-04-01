package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
)

type Portal struct {
	// APIルートURL
	API   string `config:"portal_api,required" json:"api"`
	Token string `config:"portal_token,required" json:"token"`
}

type Auth0 struct {
	Domain       string `config:"auth0_domain,required" json:"domain"`
	ClientID     string `config:"auth0_client_id,required" json:"client_id"`
	ClientSecret string `config:"auth0_client_secret,required" json:"client_secret"`
}

type Config struct {
	Portal Portal `json:"portal"`
	Auth0  Auth0  `json:"auth0"`
}

func ReadConfig(ctx context.Context) (*Config, error) {
	loader := confita.NewLoader(
		env.NewBackend(),
	)

	config := Config{}

	if err := loader.Load(ctx, &config); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	portalAPIURL, err := url.Parse(config.Portal.API)

	if err != nil {
		return nil, fmt.Errorf("envvar \"PORTAL_API\" must be a valid URL that indicates API Root of MIS.W Portal: %w", err)
	}

	if !strings.HasSuffix(portalAPIURL.Path, "/") {
		log.Println("there is no \"/\" at the end of \"PORTAL_API\". you may have problems around relative paths calculation.")
	}

	return &config, nil
}
