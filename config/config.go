package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

type Portal struct {
	// APIルートURL
	API   string
	Token string
}

type Auth0 struct {
	Domain       string
	ClientID     string
	ClientSecret string
}

type Config struct {
	Portal Portal
	Auth0  Auth0
}

func ReadConfig() (*Config, error) {
	config := Config{}

	portalAPI := os.Getenv("PORTAL_API")

	portalAPIURL, err := url.Parse(portalAPI)

	if err != nil {
		return nil, fmt.Errorf("envvar \"PORTAL_API\" must be a valid URL that indicates API Root of MIS.W Portal: %w", err)
	}

	if !strings.HasSuffix(portalAPIURL.Path, "/") {
		log.Println("there is no \"/\" at the end of \"PORTAL_API\". you may have problems around relative paths calculation.")
	}

	config.Portal.API = portalAPI

	portalToken := os.Getenv("PORTAL_TOKEN")

	if portalToken == "" {
		return nil, fmt.Errorf("envvar \"PORTAL_TOKEN\" is required")
	}

	config.Portal.Token = portalToken

	auth0Domain := os.Getenv("AUTH0_DOMAIN")

	if auth0Domain == "" {
		return nil, fmt.Errorf("envvar \"AUTH0_DOMAIN\" is required")
	}

	config.Auth0.Domain = auth0Domain

	auth0ClientID := os.Getenv("AUTH0_CLIENT_ID")

	if auth0ClientID == "" {
		return nil, fmt.Errorf("envvar \"AUTH0_CLIENT_ID\" is required")
	}

	config.Auth0.ClientID = auth0ClientID

	auth0ClientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	if auth0ClientSecret == "" {
		return nil, fmt.Errorf("envvar \"AUTH0_CLIENT_SECRET\" is required")
	}

	config.Auth0.ClientSecret = auth0ClientSecret

	return &config, nil
}
