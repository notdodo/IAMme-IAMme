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
	logger := logging.NewLogManager()
	ctx, client, err := okta.NewClient(context.Background(), okta.WithOrgUrl(fmt.Sprintf("https://%s", orgUrl)), okta.WithToken(apiKey))
	if err != nil {
		logger.Error("Invalid Okta login", "err", err)
	}
	return &oktaClient{
		oktaClient: client,
		context:    ctx,
		log:        logger,
	}
}

func (c *oktaClient) GetUsers() ([]*okta.User, error) {
	users, _, err := c.oktaClient.User.ListUsers(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *oktaClient) GetGroups() ([]*okta.Group, error) {
	groups, _, err := c.oktaClient.Group.ListGroups(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return groups, nil
}
