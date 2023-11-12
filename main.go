package main

import (
	"log"

	"github.com/notdodo/IAMme-IAMme/cmd"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"

	"github.com/joho/godotenv"
)

func main() {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		log.Fatalln(err.Error())
	}
	cmd.Execute(
		okta.NewOktaClient(envFile["OKTA_CLIENT_ORGURL"], envFile["OKTA_CLIENT_TOKEN"]),
		neo4j.NewNeo4jClient(envFile["NEO4J_URL"], envFile["NEO4J_USER"], envFile["NEO4J_PASS"]))
}
