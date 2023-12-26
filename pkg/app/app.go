package app

import (
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"
	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"
	"github.com/sourcegraph/conc/iter"

	"github.com/notdodo/goflat/v2"
	oktaSdk "github.com/okta/okta-sdk-golang/v2/okta"
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
	a.createNodes([]string{"User"}, flat(a.users()))
	groups := a.groupsWithMembers()
	a.createNodes([]string{"Group"}, flat(groups))
	rules := a.rules()
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

	a.createNodes([]string{"Application"}, flat(a.applications()))
}

func (a *iamme) users() []*User {
	oktaUsers, err := a.oktaClient.Users()
	users := make([]*User, 0, len(oktaUsers))
	if err != nil {
		a.logger.Error("Error fetching users from Okta:", "err", err)
	}

	for _, user := range oktaUsers {
		users = append(users, &User{
			Id:   user.Id,
			User: user,
		})
	}
	return users
}

func (a *iamme) groupsWithMembers() []*Group {
	oktaGroups, err := a.oktaClient.Groups()
	if err != nil {
		a.logger.Error("Error fetching groups from Okta:", "err", err)
		return nil
	}
	groupsWithMembers := iter.Map(oktaGroups, func(group **oktaSdk.Group) *Group {
		members, err := a.oktaClient.GroupMembers((*group).Id)
		if err != nil {
			a.logger.Error("Error fetching group members from Okta:", "err", err)
		}
		users := make([]*User, 0, len(members))
		for _, member := range members {
			users = append(users, &User{
				Id:   member.Id,
				User: member,
			})
		}
		return &Group{
			Id:      (*group).Id,
			Group:   *group,
			Members: users,
		}
	})

	return groupsWithMembers
}

func (a *iamme) rules() []*GroupRule {
	oktaRules, err := a.oktaClient.GroupsRules()
	rules := make([]*GroupRule, 0, len(oktaRules))
	if err != nil {
		a.logger.Error("Error fetching rules from Okta:", "err", err)
	}

	for _, rule := range oktaRules {
		rules = append(rules, &GroupRule{
			Id:        rule.Id,
			GroupRule: rule,
		})
	}
	return rules
}

func (a *iamme) applications() []*Application {
	oktaApps, err := a.oktaClient.Applications()
	apps := make([]*Application, 0, len(oktaApps))
	if err != nil {
		a.logger.Error("Error fetching rules from Okta:", "err", err)
	}

	for _, app := range oktaApps {
		apps = append(apps, &Application{
			Id:          app.(*oktaSdk.Application).Id,
			Application: app.(*oktaSdk.Application),
		})
	}
	return apps
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
