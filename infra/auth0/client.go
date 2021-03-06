//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE

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
	domain      string
	tokenSource oauth2.TokenSource
	httpClient  *http.Client
}

var _ Client = &client{}

func NewClient(domain string, tokenSource oauth2.TokenSource) Client {
	return &client{
		domain:      domain,
		tokenSource: tokenSource,
		httpClient:  http.DefaultClient,
	}
}

type auth0ErrorJSON struct {
	Message string `json:"message"`
}

// UpdateRuleConfig auth0のmanagement-apiを使ってメンバー(のロールの)情報をauth0のrules configにセットする。
// (e.g. key: "members", value: ${JSON data}
// https://auth0.com/docs/api/management/v2#!/Rules_Configs/put_rules_configs_by_key
func (c *client) UpdateRuleConfig(ctx context.Context, key, value string) error {
	token, err := c.tokenSource.Token()

	if err != nil {
		return fmt.Errorf("failed to obtain token: %w", err)
	}

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
	token.SetAuthHeader(req)

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
