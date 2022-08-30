package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/MISW/portal-role-sync/config"
	"github.com/MISW/portal-role-sync/infra/auth0"
	"github.com/MISW/portal-role-sync/infra/portal"
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

	auth0MemberPortalRoles, err := auth0Client.GetUserPortalRoles(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	portalMemberRoles, err := portalClient.GetAllMemberRoles(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	for k, v := range portalMemberRoles {
		auth0MemberPortalRole, ok := auth0MemberPortalRoles[k]
		if !ok {
			//auth0にユーザが存在しない場合、そもそもログインできないのでロール設定は実行しない。
			//ただし、会員だけどユーザが存在しないという場合に限り、createしてしまう.
			log.Printf("user {%s} with role {%s} does not exist in auth0 users.\n", k, v.Role)
			continue
		}
		if v.Role != auth0MemberPortalRole.Role {
			if err := auth0Client.UpdateUserPortalRole(ctx, k, v.Role); err != nil {
				log.Fatalln(err)
			}
		}
		log.Printf("updated user {%s} role to {%s}", k, v.Role)
	}
}
