package main

import (
	"context"
	"encoding/json"
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

	config, err := config.ReadConfig()

	if err != nil {
		log.Fatalln(err)
	}

	portalClient := portal.NewClient(config.Portal.API, config.Portal.Token)

	memberRoles, err := portalClient.GetAllMemberRoles(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	members := map[string]string{}

	for k, member := range memberRoles {
		if member.Role == "member" || member.Role == "admin" {
			members[k] = member.Role
		}
	}

	membersValue, err := json.Marshal(members)

	if err != nil {
		log.Fatalln(err)
	}

	authConfig := auth0.NewConfig(config.Auth0.Domain, config.Auth0.ClientID, config.Auth0.ClientSecret)

	token, err := authConfig.Token(ctx)

	if err != nil {
		panic(err)
	}

	auth0Client := auth0.NewClient(config.Auth0.Domain, token)

	if err := auth0Client.UpdateRuleConfig(ctx, "members", string(membersValue)); err != nil {
		panic(err)
	}

}
