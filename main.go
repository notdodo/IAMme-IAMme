package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/notdodo/IAMme-IAMme/cmd"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
)

func main() {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		log.Fatalln(err.Error())
	}
	cmd.Execute(neo4j.NewNeo4jClient(envFile["NEO4J_URL"], envFile["NEO4J_USER"], envFile["NEO4J_PASS"]))
}
