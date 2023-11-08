package app

import (
	"context"
	"log"
	"strings"

	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"

	neo4jSdk "github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
		log.Println(err.Error())
	}

	userParams := make([]map[string]interface{}, 0)
	for _, user := range users {
		userParams = append(userParams, map[string]interface{}{
			"userId":    user.Id,
			"status":    user.Status,
			"firstName": (*user.Profile)["firstName"],
			"lastName":  (*user.Profile)["lastName"],
		})
	}

	session := a.neo4jClient.Connect()
	ctx := context.TODO()
	query := buildDynamicQuery(userParams)
	_, err = session.ExecuteWrite(ctx, func(tx neo4jSdk.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{
			"userParams": userParams,
		})

		if err != nil {
			panic(err)
		}
		return nil, err
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// TODO: this is ugly AF
func buildDynamicQuery(userParams []map[string]interface{}) string {
	var queryBuilder strings.Builder

	queryBuilder.WriteString("UNWIND $userParams as user\n")
	queryBuilder.WriteString("CREATE (u:User {\n")

	fields := userParams[0]

	i := 0
	for field := range fields {
		queryBuilder.WriteString(fieldKeyToCypherProperty(field) + ": user." + fieldKeyToCypherProperty(field))
		i++
		if i < len(fields) {
			queryBuilder.WriteString(",\n")
		}
	}

	queryBuilder.WriteString("\n})\n")
	queryBuilder.WriteString("RETURN u\n")

	return queryBuilder.String()
}

// TODO: this is ugly AF
func fieldKeyToCypherProperty(key interface{}) string {
	keyStr, ok := key.(string)
	if !ok {
		panic("Invalid field key")
	}

	return keyStr
}
