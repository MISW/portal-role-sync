package auth0

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

// Auth0のクライアント
type Client interface {
	UpdateRuleConfig(ctx context.Context, key, value string) error
}

type client struct {
	domain     string
	token      *oauth2.Token
	httpClient *http.Client
}

var _ Client = &client{}

func NewClient(domain string, token *oauth2.Token) Client {
	return &client{
		domain:     domain,
		token:      token,
		httpClient: http.DefaultClient,
	}
}

type auth0ErrorJSON struct {
	Message string `json:"message"`
}

func (c *client) UpdateRuleConfig(ctx context.Context, key, value string) error {
	reqBody, err := json.Marshal(map[string]string{
		"value": value,
	})

	if err != nil {
		return fmt.Errorf("failed to marshal rule config update request payload: %w", err)
	}

	endpoint := "https://" + c.domain + "/api/v2/rules-configs/" + key

	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewBuffer(reqBody))

	if err != nil {
		return fmt.Errorf("failed to create rule config update request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.token.SetAuthHeader(req)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send rule config update request: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var auth0Error auth0ErrorJSON
		if err := json.Unmarshal(respBody, &auth0Error); err != nil {
			return fmt.Errorf("failed to parse auth0 error message: %w", err)
		}
		return fmt.Errorf("failed to update role config: %s %s", resp.Status, auth0Error.Message)
	}

	return nil
}
