package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/MISW/portal-role-sync/config"
	"github.com/MISW/portal-role-sync/infra/auth0"
	"github.com/MISW/portal-role-sync/infra/portal"
	"github.com/MISW/portal-role-sync/reconciler"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config, err := config.ReadConfig(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	portalClient := portal.NewClient(config.Portal.API, config.Portal.Token)

	auth0AuthConfig := auth0.NewConfig(config.Auth0.Domain, config.Auth0.ClientID, config.Auth0.ClientSecret)
	auth0Client := auth0.NewClient(config.Auth0.Domain, auth0AuthConfig.TokenSource(ctx))

	auth0Reconciler := reconciler.NewAuth0Reconciler(auth0Client)

	memberRoles, err := portalClient.GetAllMemberRoles(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	req := &reconciler.ReconcileRequest{
		Members: memberRoles,
	}

	if err := auth0Reconciler.Reconcile(ctx, req); err != nil {
		log.Fatalln(err)
	}
}
