package app

import (
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"
	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"

	"github.com/notdodo/goflat"
)

type IAMme interface {
	Dump()
}

func NewIAMme(oktaClient okta.OktaClient, neo4jClient neo4j.Neo4jClient) IAMme {
	logger := logging.GetLogManager()
	return &iamme{
		oktaClient:  oktaClient,
		neo4jClient: neo4jClient,
		logger:      logger,
	}
}

type iamme struct {
	oktaClient  okta.OktaClient
	neo4jClient neo4j.Neo4jClient
	logger      logging.LogManager
}

func (a *iamme) Dump() {
	a.createNodes([]string{"User"}, flat(a.getUsers()))
	a.createNodes([]string{"Group"}, flat(a.getGroups()))
	rules := a.getRules()
	a.createNodes([]string{"Rule"}, flat(rules))

	groupRules := make([]map[string]interface{}, 0)
	for _, rule := range rules {
		for _, gid := range rule.Actions.AssignUserToGroups.GroupIds {
			groupRules = append(groupRules, map[string]interface{}{
				"left_key":    "GroupRule_Id",
				"left_value":  rule.GroupRule.Id,
				"right_key":   "Group_Id",
				"right_value": gid,
			})
		}
	}
	a.createRelations([]string{"GroupRule"}, []string{"Rule"}, []string{"Group"}, groupRules)
}

func (a *iamme) getUsers() []*okta.User {
	users, err := a.oktaClient.GetUsers()
	if err != nil {
		a.logger.Error("Error fetching users from Okta:", "err", err)
	}
	return users
}

func (a *iamme) getGroups() []*okta.Group {
	groups, err := a.oktaClient.GetGroups()
	if err != nil {
		a.logger.Error("Error fetching groups from Okta:", "err", err)
	}
	return groups
}

func (a *iamme) getRules() []*okta.GroupRule {
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

func (a *iamme) createNodes(nodeLabels []string, properties []map[string]interface{}) []map[string]interface{} {
	nodeIDs, err := a.neo4jClient.CreateNodes(nodeLabels, properties)
	if err != nil {
		a.logger.Error("Error creating nodes on Neo4J", "err", err)
	}
	return nodeIDs
}

func (a *iamme) createRelations(relationLabels []string, aLabels []string, bLabels []string, properties []map[string]interface{}) []map[string]interface{} {
	relIDs, err := a.neo4jClient.CreateRelationsAtoB(relationLabels, aLabels, bLabels, properties)
	if err != nil {
		a.logger.Error("Error creating nodes on Neo4J", "err", err)
	}
	return relIDs
}
