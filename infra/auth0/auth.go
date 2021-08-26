package auth0

import (
	"net/url"

	"golang.org/x/oauth2/clientcredentials"
)

// Auth0のManagement APIを叩くためのトークンを取得するConfigを生成する
func NewConfig(domain, clientID, clientSecret string) *clientcredentials.Config {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://" + domain + "/oauth/token",
		EndpointParams: url.Values{
			"audience": {"https://" + domain + "/api/v2/"},
		},
	}

	return &config
}
