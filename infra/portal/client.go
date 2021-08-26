package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MemberRole struct {
	Role string `json:"role"`
}

type MemberRoles = map[string]MemberRole

// みすポータルのクライアント
type Client interface {
	GetAllMemberRoles(ctx context.Context) (MemberRoles, error)
}

type client struct {
	root       string
	token      string
	httpClient *http.Client
}

var _ Client = &client{}

func NewClient(apiRoot string, token string) Client {
	return &client{
		root:       apiRoot,
		token:      token,
		httpClient: http.DefaultClient,
	}
}

type portalError struct {
	Message string `json:"message"`
}

func (c *client) GetAllMemberRoles(ctx context.Context) (MemberRoles, error) {
	endpoint := c.root + "external/all_member_roles"

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errorMessage portalError
		if err = json.Unmarshal(body, &errorMessage); err != nil {
			return nil, fmt.Errorf("failed to fetch member roles (status: %d)", resp.StatusCode)
		}
		return nil, fmt.Errorf("failed to fetch member roles (status: %d, message: \"%s\")", resp.StatusCode, errorMessage.Message)
	}

	var memberRoles MemberRoles
	if err = json.Unmarshal(body, &memberRoles); err != nil {
		return nil, err
	}

	return memberRoles, nil
}
