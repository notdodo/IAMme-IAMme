package okta

import (
	"context"
	"fmt"

	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

// OktaClient is an interface for interacting with Okta resources.
type OktaClient interface {
	GetUsers() ([]*User, error)
	GetGroups() ([]*Group, error)
	GetGroupsRules() ([]*GroupRule, error)
	GetGroupMembers(string) ([]*User, error)
}

type oktaClient struct {
	oktaClient *okta.Client
	log        logging.LogManager
}

type Group struct {
	*okta.Group
	Members []*User
}

type User struct {
	*okta.User
}

type GroupRule struct {
	*okta.GroupRule
}

func NewOktaClient(orgUrl, apiKey string) OktaClient {
	logger := logging.GetLogManager()
	_, client, err := okta.NewClient(context.TODO(), okta.WithOrgUrl(fmt.Sprintf("https://%s", orgUrl)), okta.WithToken(apiKey))
	if err != nil {
		logger.Error("Invalid Okta login", "err", err)
	}
	logger.Info("Valid Okta client", "orgUrl", orgUrl)
	return &oktaClient{
		oktaClient: client,
		log:        logger,
	}
}

func (c *oktaClient) GetUsers() ([]*User, error) {
	c.log.Info("Getting Okta users")
	var users []*User
	appendUsers := func(oktaUsers []*okta.User) {
		for _, oktaUser := range oktaUsers {
			users = append(users, &User{
				User: oktaUser,
			})
		}
	}

	oktaUsers, response, err := c.oktaClient.User.ListUsers(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	appendUsers(oktaUsers)
	for response.HasNextPage() {
		response, err = response.Next(context.TODO(), &oktaUsers)
		if err != nil {
			return nil, err
		}
		appendUsers(oktaUsers)
	}

	c.log.Info(fmt.Sprintf("Found %d users", len(users)))
	c.log.Debug(fmt.Sprintf("Found %d users", len(users)), "users", response.Body)
	return users, nil
}

func (c *oktaClient) GetGroups() ([]*Group, error) {
	c.log.Info("Getting Okta groups")
	oktaGroups, response, err := c.oktaClient.Group.ListGroups(context.TODO(), &query.Params{
		Expand: "stats,app",
	})
	if err != nil {
		return nil, err
	}
	groups := make([]*Group, 0, len(oktaGroups))
	for _, oktaGroup := range oktaGroups {
		groups = append(groups, &Group{
			Group: oktaGroup,
		})
	}

	c.log.Info(fmt.Sprintf("Found %d groups", len(groups)))
	c.log.Debug(fmt.Sprintf("Found %d groups", len(groups)), "groups", response.Body)
	return groups, nil
}

func (c *oktaClient) GetGroupsRules() ([]*GroupRule, error) {
	c.log.Info("Getting Okta groups rules")
	oktaRules, response, err := c.oktaClient.Group.ListGroupRules(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	rules := make([]*GroupRule, 0, len(oktaRules))
	for _, oktaRule := range oktaRules {
		rules = append(rules, &GroupRule{
			GroupRule: oktaRule,
		})
	}
	c.log.Info(fmt.Sprintf("Found %d rules", len(rules)))
	c.log.Debug(fmt.Sprintf("Found %d rules", len(rules)), "rules", response.Body)
	return rules, nil
}

func (c *oktaClient) GetGroupMembers(groupId string) ([]*User, error) {
	c.log.Info("Getting Okta group members", "group", groupId)
	oktaMembers, response, err := c.oktaClient.Group.ListGroupUsers(context.TODO(), groupId, nil)
	if err != nil {
		return nil, err
	}
	members := make([]*User, 0, len(oktaMembers))
	for _, member := range oktaMembers {
		members = append(members, &User{
			User: member,
		})
	}
	c.log.Info(fmt.Sprintf("Found %d members for group %s", len(members), groupId))
	c.log.Debug(fmt.Sprintf("Found %d members for group %s", len(members), groupId), "members", response.Body)
	return members, err
}
