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
	logger := logging.GetLogManager()
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
	a.fetchAndCreateNodes([]string{"User"}, a.getUsers)
	a.fetchAndCreateNodes([]string{"Groups"}, a.getGroups)
}

func (a *oktaNeo4jApp) getUsers() ([]interface{}, error) {
	users, err := a.oktaClient.GetUsers()
	if err != nil {
		a.logger.Error("Error fetching users from Okta:", "err", err)
		return nil, err
	}
	usersInterface := make([]interface{}, len(users))
	for i, user := range users {
		usersInterface[i] = user
	}

	return usersInterface, nil
}

func (a *oktaNeo4jApp) getGroups() ([]interface{}, error) {
	groups, err := a.oktaClient.GetGroups()
	if err != nil {
		a.logger.Error("Error fetching groups from Okta:", "err", err)
		return nil, err
	}
	groupsInterface := make([]interface{}, len(groups))
	for i, group := range groups {
		groupsInterface[i] = group
	}

	return groupsInterface, nil
}

func (a *oktaNeo4jApp) fetchAndCreateNodes(nodeLabels []string, oktaClientFunc func() ([]interface{}, error)) {
	data, err := oktaClientFunc()
	if err != nil {
		a.logger.Error("Error fetching data from Okta:", "err", err)
		return
	}

	flatData := make([]map[string]interface{}, 0, len(data))
	for _, item := range data {
		flatItem := goflat.FlatStruct(item, goflat.FlattenerConfig{
			Separator: "_",
			OmitEmpty: true,
			OmitNil:   true,
		})
		flatData = append(flatData, flatItem)
	}

	if _, err := a.neo4jClient.CreateNodes(nodeLabels, &flatData); err != nil {
		a.logger.Error("Error creating nodes on Neo4J", "err", err)
	}
}
