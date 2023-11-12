package app

import (
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"
	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"

	"github.com/notdodo/goflat"
)

type OktaNeo4jApp interface {
	Dump()
}

func NewOktaNeo4jApp(oktaClient okta.OktaClient, neo4jClient neo4j.Neo4jClient) OktaNeo4jApp {
	logger := logging.NewLogManager()
	return &oktaNeo4jApp{
		oktaClient:  oktaClient,
		neo4jClient: neo4jClient,
		logger:      logger,
	}
}

type oktaNeo4jApp struct {
	oktaClient  okta.OktaClient
	neo4jClient neo4j.Neo4jClient
	logger      logging.LogManager
}

func (a *oktaNeo4jApp) Dump() {
	users, err := a.oktaClient.GetUsers()
	if err != nil {
		a.logger.Error("Error fetching users from Okta:", "err", err)
		return
	}

	flatUsers := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		flatUser := goflat.FlatStruct(*user, goflat.FlattenerConfig{
			Separator: "_",
			OmitEmpty: true,
			OmitNil:   true,
		})
		flatUsers = append(flatUsers, flatUser)
	}

	if _, err = a.neo4jClient.CreateNodes([]string{"User"}, &flatUsers); err != nil {
		a.logger.Error("Error creating user nodes on Neo4J", "err", err)
	}
}
