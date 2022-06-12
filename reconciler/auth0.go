package reconciler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MISW/portal-role-sync/infra/auth0"
)

func NewAuth0Reconciler(auth0Client auth0.Client) Reconciler {
	return &auth0Reconciler{
		auth0Client: auth0Client,
	}
}

type auth0Reconciler struct {
	auth0Client auth0.Client
}

func (r *auth0Reconciler) Reconcile(ctx context.Context, req *ReconcileRequest) error {
	members := map[string]string{}

	for key, member := range req.Members {
		if member.Role == "member" || member.Role == "admin" {
			members[key] = member.Role
		}
	}

	membersJSON, err := json.Marshal(members)
	log.Print(members)

	if err != nil {
		return fmt.Errorf("failed to marshal members: %w", err)
	}

	if r.auth0Client.UpdateRuleConfig(ctx, "members", string(membersJSON)); err != nil {
		return fmt.Errorf("failed to update rule config: %w", err)
	}

	return nil
}
