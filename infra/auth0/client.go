//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../mock/$GOPACKAGE/$GOFILE

package auth0

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MISW/portal-role-sync/infra/portal"
	"golang.org/x/oauth2"
)

// Auth0のクライアント
type Client interface {
	UpdateRuleConfig(ctx context.Context, key, value string) error
	GetUserPortalRoles(ctx context.Context) (portal.MemberRoles, error)
	UpdateUserPortalRole(ctx context.Context, roleKey, role string) error
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
type auth0AppMetadata struct {
	PortalRole string `json:"portal_role"`
}
type auth0Member struct {
	UserID          string           `json:"user_id"`
	AppMetadataRole auth0AppMetadata `json:"app_metadata"`
}

const userIdSlackPrefix = "oauth2|slack|"

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

// GetUserPortalRoles auth0のユーザーのapp_metadataを見てportal_roleを取得する
// 参考: https://auth0.com/docs/manage-users/user-search/retrieve-users-with-get-users-endpoint
func (c *client) GetUserPortalRoles(ctx context.Context) (portal.MemberRoles, error) {
	token, err := c.tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to obtain token: %w", err)
	}

	endpoint := fmt.Sprintf("https://%s/api/v2/users", c.domain)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user app_metadata update request: %w", err)
	}
	token.SetAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send rule config update request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var auth0Error auth0ErrorJSON
		if err := json.Unmarshal(respBody, &auth0Error); err != nil {
			return nil, fmt.Errorf("failed to parse auth0 error message: %w", err)
		}
		return nil, fmt.Errorf("failed to get users: %s %s", resp.Status, auth0Error.Message)
	}

	//fmt.Println(string(respBody))
	var auth0Members []auth0Member
	if err := json.Unmarshal(respBody, &auth0Members); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	portalMemberRoles := make(portal.MemberRoles, 0)
	for _, v := range auth0Members {
		if !strings.HasPrefix(v.UserID, userIdSlackPrefix) {
			return nil, fmt.Errorf("failed to find slack_id prefix in user_od. %s doesn's have prefix %s", v.UserID, userIdSlackPrefix)
		}
		slackID := strings.TrimPrefix(v.UserID, userIdSlackPrefix)
		portalMemberRoles[slackID] = portal.MemberRole{
			Role: v.AppMetadataRole.PortalRole,
		}
	}

	return portalMemberRoles, nil
}

//UpdateUserPortalRoles appMetaDataにみすポータルのロール情報をセットする
func (c *client) UpdateUserPortalRole(ctx context.Context, slackID, role string) error {
	token, err := c.tokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to obtain token: %w", err)
	}

	data := map[string]auth0AppMetadata{
		"app_metadata": {
			PortalRole: role, //nullに設定したら削除される。
		},
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal json. %w", err)
	}

	userID := fmt.Sprintf("%s%s", userIdSlackPrefix, slackID)
	endpoint := fmt.Sprintf("https://%s/api/v2/users/%s", c.domain, userID)
	req, err := http.NewRequestWithContext(ctx, "PATCH", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create user app_metadata update request: %w", err)
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
		return fmt.Errorf("failed to update user %s's app_metadata portal_role: %s %s", userID, resp.Status, auth0Error.Message)
	}
	return nil
}
