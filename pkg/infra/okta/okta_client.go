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
	ctx        context.Context
}

func NewOktaClient(orgUrl, apiKey string) OktaClient {
	logger := logging.GetLogManager()
	ctx, client, err := okta.NewClient(context.TODO(), okta.WithOrgUrl(fmt.Sprintf("https://%s", orgUrl)), okta.WithToken(apiKey))
	if err != nil {
		logger.Error("Invalid Okta login", "err", err)
	}
	logger.Info("Valid Okta client", "org_url", orgUrl)
	return &oktaClient{
		oktaClient: client,
		log:        logger,
		ctx:        ctx,
	}
}

func (c *oktaClient) Users() ([]*okta.User, error) {
	c.log.Info("Getting Okta users")
	var users []*okta.User

	oktaUsers, response, err := c.oktaClient.User.ListUsers(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	users = append(users, oktaUsers...)
	for response.HasNextPage() {
		var oktaUsers []*okta.User
		response, err = response.Next(context.TODO(), &oktaUsers)
		if err != nil {
			return nil, err
		}
		users = append(users, oktaUsers...)
	}

	c.log.Info(fmt.Sprintf("Found %d users", len(users)))
	c.log.Debug(fmt.Sprintf("Found %d users", len(users)), "users", users)
	return users, nil
}

func (c *oktaClient) Groups() ([]*okta.Group, error) {
	c.log.Info("Getting Okta groups")
	oktaGroups, _, err := c.oktaClient.Group.ListGroups(context.TODO(), &query.Params{
		Expand: "stats,app",
		Limit:  10000,
	})
	if err != nil {
		return nil, err
	}

	c.log.Info(fmt.Sprintf("Found %d groups", len(oktaGroups)))
	c.log.Debug(fmt.Sprintf("Found %d groups", len(oktaGroups)), "groups", oktaGroups)
	return oktaGroups, nil
}

func (c *oktaClient) GroupsRules() ([]*okta.GroupRule, error) {
	c.log.Info("Getting Okta groups rules")
	var rules []*okta.GroupRule

	oktaRules, response, err := c.oktaClient.Group.ListGroupRules(context.TODO(), &query.Params{
		Limit: 200,
	})
	if err != nil {
		return nil, err
	}
	rules = append(rules, oktaRules...)
	for response.HasNextPage() {
		var oktaRules []*okta.GroupRule
		response, err = response.Next(context.TODO(), &oktaRules)
		if err != nil {
			return nil, err
		}
		rules = append(rules, oktaRules...)
	}

	c.log.Info(fmt.Sprintf("Found %d rules", len(rules)))
	c.log.Debug(fmt.Sprintf("Found %d rules", len(rules)), "rules", rules)
	return rules, nil
}

func (c *oktaClient) GroupMembers(groupId string) ([]*okta.User, error) {
	c.log.Info("Getting Okta group members", "group", groupId)
	var members []*okta.User

	oktaMembers, response, err := c.oktaClient.Group.ListGroupUsers(context.TODO(), groupId, &query.Params{
		Limit: 1000,
	})
	if err != nil {
		return nil, err
	}
	members = append(members, oktaMembers...)
	for response.HasNextPage() {
		var oktaMembers []*okta.User
		response, err = response.Next(context.TODO(), &oktaMembers)
		if err != nil {
			return nil, err
		}
		members = append(members, oktaMembers...)
	}

	c.log.Info(fmt.Sprintf("Found %d members for group %s", len(members), groupId))
	c.log.Debug(fmt.Sprintf("Found %d members for group %s", len(members), groupId), "members", members)
	return oktaMembers, err
}

func (c *oktaClient) Applications() ([]okta.App, error) {
	c.log.Info("Getting Okta applications")
	var apps []okta.App

	oktaApps, response, err := c.oktaClient.Application.ListApplications(context.TODO(), &query.Params{
		Limit: 200,
	})
	if err != nil {
		return nil, err
	}
	apps = append(apps, oktaApps...)
	for response.HasNextPage() {
		var oktaApps []okta.App
		response, err = response.Next(context.TODO(), &oktaApps)
		if err != nil {
			return nil, err
		}
		apps = append(apps, oktaApps...)
	}

	c.log.Info(fmt.Sprintf("Found %d applications", len(apps)))
	c.log.Debug(fmt.Sprintf("Found %d applications", len(apps)), "applications", apps)
	return oktaApps, err
}

func (c *oktaClient) ApplicationGroupAssignments(appId string) ([]*okta.ApplicationGroupAssignment, error) {
	c.log.Info("Getting Okta application group assigments", "application", appId)
	oktaGroups, _, err := c.oktaClient.Application.ListApplicationGroupAssignments(context.TODO(), appId, &query.Params{
		Limit: 200,
	})
	if err != nil {
		return nil, err
	}
	c.log.Info(fmt.Sprintf("Found %d group assigments for application %s", len(oktaGroups), appId))
	c.log.Debug(fmt.Sprintf("Found %d group assigments for application %s", len(oktaGroups), appId), "groupassignments", oktaGroups)
	return oktaGroups, err
}
