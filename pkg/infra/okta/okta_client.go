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
	GetUsers() ([]*okta.User, error)
	GetGroups() ([]*okta.Group, error)
	GetGroupsRules() ([]*okta.GroupRule, error)
	GetGroupMembers(string) ([]*okta.User, error)
}

type oktaClient struct {
	oktaClient *okta.Client
	log        logging.LogManager
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

func (c *oktaClient) GetUsers() ([]*okta.User, error) {
	c.log.Info("Getting Okta users")
	var users []*okta.User
	appendUsers := func(oktaUsers []*okta.User) {
		for _, oktaUser := range oktaUsers {
			users = append(users, oktaUser)
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

func (c *oktaClient) GetGroups() ([]*okta.Group, error) {
	c.log.Info("Getting Okta groups")
	oktaGroups, response, err := c.oktaClient.Group.ListGroups(context.TODO(), &query.Params{
		Expand: "stats,app",
	})
	if err != nil {
		return nil, err
	}

	c.log.Info(fmt.Sprintf("Found %d groups", len(oktaGroups)))
	c.log.Debug(fmt.Sprintf("Found %d groups", len(oktaGroups)), "groups", response.Body)
	return oktaGroups, nil
}

func (c *oktaClient) GetGroupsRules() ([]*okta.GroupRule, error) {
	c.log.Info("Getting Okta groups rules")
	oktaRules, response, err := c.oktaClient.Group.ListGroupRules(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	c.log.Info(fmt.Sprintf("Found %d rules", len(oktaRules)))
	c.log.Debug(fmt.Sprintf("Found %d rules", len(oktaRules)), "rules", response.Body)
	return oktaRules, nil
}

func (c *oktaClient) GetGroupMembers(groupId string) ([]*okta.User, error) {
	c.log.Info("Getting Okta group members", "group", groupId)
	oktaMembers, response, err := c.oktaClient.Group.ListGroupUsers(context.TODO(), groupId, nil)
	if err != nil {
		return nil, err
	}
	c.log.Info(fmt.Sprintf("Found %d members for group %s", len(oktaMembers), groupId))
	c.log.Debug(fmt.Sprintf("Found %d members for group %s", len(oktaMembers), groupId), "members", response.Body)
	return oktaMembers, err
}
