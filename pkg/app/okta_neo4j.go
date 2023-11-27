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
	a.createNodes([]string{"User"}, flat(a.getUsers()))
	a.createNodes([]string{"Group"}, flat(a.getGroups()))
	a.createNodes([]string{"Rule"}, flat(a.getRules()))

	// var groupIdsList [][]string

	// for _, rule := range rules {
	// 	groupIds := rule.Actions.AssignUserToGroups.GroupIds
	// 	groupIdsList = append(groupIdsList, groupIds)
	// }
}

func (a *oktaNeo4jApp) getUsers() []*okta.User {
	users, err := a.oktaClient.GetUsers()
	if err != nil {
		a.logger.Error("Error fetching users from Okta:", "err", err)
	}
	return users
}

func (a *oktaNeo4jApp) getGroups() []*okta.Group {
	groups, err := a.oktaClient.GetGroups()
	if err != nil {
		a.logger.Error("Error fetching groups from Okta:", "err", err)
	}
	return groups
}

func (a *oktaNeo4jApp) getRules() []*okta.GroupRule {
	rules, err := a.oktaClient.GetGroupsRules()
	if err != nil {
		a.logger.Error("Error fetching rules from Okta:", "err", err)
	}
	return rules
}

func flat[T any](data []T) []map[string]interface{} {
	flatData := make([]map[string]interface{}, 0, len(data))
	for _, item := range data {
		flatItem := goflat.FlatStruct(item, goflat.FlattenerConfig{
			Separator: "_",
			OmitEmpty: true,
			OmitNil:   true,
		})
		flatData = append(flatData, flatItem)
	}
	return flatData
}

func (a *oktaNeo4jApp) createNodes(nodeLabels []string, properties []map[string]interface{}) ([]map[string]interface{}, error) {
	nodeIDs, err := a.neo4jClient.CreateNodes(nodeLabels, properties)
	if err != nil {
		a.logger.Error("Error creating nodes on Neo4J", "err", err)
	}
	return nodeIDs, err
}
