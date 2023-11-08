package okta

import (
	"context"
	"fmt"
	"log"

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
}

func NewOktaClient(orgUrl, apiKey string) OktaClient {
	ctx, client, err := okta.NewClient(context.Background(), okta.WithOrgUrl(fmt.Sprintf("https://%s", orgUrl)), okta.WithToken(apiKey))
	if err != nil {
		log.Fatalln("Invalid Okta login", err.Error())
	}
	return &oktaClient{
		oktaClient: client,
		context:    ctx,
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
