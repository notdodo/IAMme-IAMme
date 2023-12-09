package app

import (
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"
	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"
	"github.com/sourcegraph/conc/iter"

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
	groups := a.getGroups()
	a.createNodes([]string{"Group"}, flat(groups))
	rules := a.getRules()
	a.createNodes([]string{"Rule"}, flat(rules))

	groupRules := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		for _, gid := range rule.Actions.AssignUserToGroups.GroupIds {
			groupRules = append(groupRules, map[string]interface{}{
				"left_key":    "GroupRule_Id",
				"left_value":  rule.Id,
				"right_key":   "Group_Id",
				"right_value": gid,
			})
		}
	}
	a.createRelations("GroupRule", []string{"Rule"}, []string{"Group"}, groupRules)

	groupMembers := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		for _, gid := range group.Members {
			groupMembers = append(groupMembers, map[string]interface{}{
				"left_key":    "User_Id",
				"left_value":  gid.Id,
				"right_key":   "Group_Id",
				"right_value": group.Id,
			})
		}
	}
	a.createRelations("GroupMember", []string{"User"}, []string{"Group"}, groupMembers)
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
		return nil
	}
	groupsWithMembers := iter.Map(groups, func(group **okta.Group) *okta.Group {
		members, err := a.oktaClient.GetGroupMembers((*group).Id)
		if err != nil {
			a.logger.Error("Error fetching group members from Okta:", "err", err)
		}
		elem := &okta.Group{
			Group:   (*group).Group,
			Members: members,
		}
		return elem
	})

	return groupsWithMembers
}

func (a *iamme) getRules() []*okta.GroupRule {
	rules, err := a.oktaClient.GetGroupsRules()
	if err != nil {
		a.logger.Error("Error fetching rules from Okta:", "err", err)
	}
	return rules
}

func flat[T any](data []T) []map[string]interface{} {
	flatData := iter.Map(data, func(item *T) map[string]interface{} {
		return goflat.FlatStruct(*item, goflat.FlattenerConfig{
			Separator: "_",
			OmitEmpty: true,
			OmitNil:   true,
		})
	})
	return flatData
}

func (a *iamme) createNodes(nodeLabels []string, properties []map[string]interface{}) []map[string]interface{} {
	nodeIDs, err := a.neo4jClient.CreateNodes(nodeLabels, properties)
	if err != nil {
		a.logger.Error("Error creating nodes on Neo4J", "err", err)
	}
	return nodeIDs
}

func (a *iamme) createRelations(relationLabel string, aLabels []string, bLabels []string, properties []map[string]interface{}) []map[string]interface{} {
	relIDs, err := a.neo4jClient.CreateRelationsAtoB(relationLabel, aLabels, bLabels, properties)
	if err != nil {
		a.logger.Error("Error creating nodes on Neo4J", "err", err)
	}
	return relIDs
}
