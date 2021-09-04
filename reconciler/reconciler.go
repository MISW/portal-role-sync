package reconciler

import (
	"context"

	"github.com/MISW/portal-role-sync/infra/portal"
)

type ReconcileRequest struct {
	Members portal.MemberRoles
}

type Reconciler interface {
	Reconcile(ctx context.Context, req *ReconcileRequest) error
}
