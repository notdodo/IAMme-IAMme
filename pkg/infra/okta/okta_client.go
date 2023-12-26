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
	Users() ([]*okta.User, error)
	Groups() ([]*okta.Group, error)
	GroupsRules() ([]*okta.GroupRule, error)
	GroupMembers(string) ([]*okta.User, error)
	Applications() ([]okta.App, error)
	ApplicationGroupAssignments(string) ([]*okta.ApplicationGroupAssignment, error)
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

func (c *oktaClient) Users() ([]*okta.User, error) {
	c.log.Info("Getting Okta users")
	var users []*okta.User
	appendUsers := func(oktaUsers []*okta.User) {
		users = append(users, oktaUsers...)
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

func (c *oktaClient) Groups() ([]*okta.Group, error) {
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

func (c *oktaClient) GroupsRules() ([]*okta.GroupRule, error) {
	c.log.Info("Getting Okta groups rules")
	oktaRules, response, err := c.oktaClient.Group.ListGroupRules(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	c.log.Info(fmt.Sprintf("Found %d rules", len(oktaRules)))
	c.log.Debug(fmt.Sprintf("Found %d rules", len(oktaRules)), "rules", response.Body)
	return oktaRules, nil
}

func (c *oktaClient) GroupMembers(groupId string) ([]*okta.User, error) {
	c.log.Info("Getting Okta group members", "group", groupId)
	oktaMembers, response, err := c.oktaClient.Group.ListGroupUsers(context.TODO(), groupId, nil)
	if err != nil {
		return nil, err
	}
	c.log.Info(fmt.Sprintf("Found %d members for group %s", len(oktaMembers), groupId))
	c.log.Debug(fmt.Sprintf("Found %d members for group %s", len(oktaMembers), groupId), "members", response.Body)
	return oktaMembers, err
}

func (c *oktaClient) Applications() ([]okta.App, error) {
	c.log.Info("Getting Okta applications")
	oktaApps, response, err := c.oktaClient.Application.ListApplications(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	c.log.Info(fmt.Sprintf("Found %d applications", len(oktaApps)))
	c.log.Debug(fmt.Sprintf("Found %d applications", len(oktaApps)), "applications", response.Body)
	return oktaApps, err
}

func (c *oktaClient) ApplicationGroupAssignments(appId string) ([]*okta.ApplicationGroupAssignment, error) {
	c.log.Info("Getting Okta application group assigments", "application", appId)
	oktaGroups, response, err := c.oktaClient.Application.ListApplicationGroupAssignments(context.TODO(), appId, nil)
	if err != nil {
		return nil, err
	}
	c.log.Info(fmt.Sprintf("Found %d group assigments for application %s", len(oktaGroups), appId))
	c.log.Debug(fmt.Sprintf("Found %d group assigments for application %s", len(oktaGroups), appId), "groupassignments", response.Body)
	return oktaGroups, err
}
