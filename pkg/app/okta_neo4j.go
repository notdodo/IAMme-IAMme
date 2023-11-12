package app

import (
	"log"

	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"
	"github.com/notdodo/goflat"
)

type OktaNeo4jApp interface {
	Dump()
}

func NewOktaNeo4jApp(oktaClient okta.OktaClient, neo4jClient neo4j.Neo4jClient) OktaNeo4jApp {
	return &oktaNeo4jApp{
		oktaClient:  oktaClient,
		neo4jClient: neo4jClient,
	}
}

type oktaNeo4jApp struct {
	oktaClient  okta.OktaClient
	neo4jClient neo4j.Neo4jClient
}

func (a *oktaNeo4jApp) Dump() {
	users, err := a.oktaClient.GetUsers()
	if err != nil {
		log.Println("Error fetching users from Okta:", err)
		return
	}

	flatUsers := make([]map[string]interface{}, 0)

	for _, user := range users {
		flatUser := goflat.FlatStruct(user, goflat.FlattenerConfig{
			Separator: "_",
			OmitEmpty: true,
			OmitNil:   true,
		})
		flatUsers = append(flatUsers, flatUser)
	}

	if _, err = a.neo4jClient.CreateNodes([]string{"User"}, &flatUsers); err != nil {
		log.Fatalln(err)
	}
}
