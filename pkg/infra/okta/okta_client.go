package okta

import (
	"context"
	"fmt"

	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// OktaClient is an interface for interacting with Okta resources.
type OktaClient interface {
	GetUsers() ([]*okta.User, error)
	GetGroups() ([]*okta.Group, error)
}

type oktaClient struct {
	oktaClient *okta.Client
	context    context.Context
	log        logging.LogManager
}

func NewOktaClient(orgUrl, apiKey string) OktaClient {
	logger := logging.GetLogManager()
	ctx, client, err := okta.NewClient(context.TODO(), okta.WithOrgUrl(fmt.Sprintf("https://%s", orgUrl)), okta.WithToken(apiKey))
	if err != nil {
		logger.Error("Invalid Okta login", "err", err)
	}
	logger.Info("Valid Okta client", "orgUrl", orgUrl)
	return &oktaClient{
		oktaClient: client,
		context:    ctx,
		log:        logger,
	}
}

func (c *oktaClient) GetUsers() ([]*okta.User, error) {
	c.log.Info("Getting Okta users")
	users, response, err := c.oktaClient.User.ListUsers(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	c.log.Info(fmt.Sprintf("Found %d users", len(users)))
	c.log.Debug(fmt.Sprintf("Found %d users", len(users)), "users", response.Body)
	return users, nil
}

func (c *oktaClient) GetGroups() ([]*okta.Group, error) {
	c.log.Info("Getting Okta groups")
	groups, response, err := c.oktaClient.Group.ListGroups(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	c.log.Info(fmt.Sprintf("Found %d groups", len(groups)))
	c.log.Debug(fmt.Sprintf("Found %d groups", len(groups)), "groups", response.Body)
	return groups, nil
}
