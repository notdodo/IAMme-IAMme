package cmd

import (
	"fmt"

	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"

	"github.com/spf13/cobra"
)

type clients_type struct {
	okta  okta.OktaClient
	neo4j neo4j.Neo4jClient
}

var clients *clients_type
var rootCmd = &cobra.Command{
	Use:   "github.com/notdodo/IAMme-IAMme",
	Short: "A CLI tool to interact with Okta and Neo4j",
}

func Execute(oktaClient okta.OktaClient, neo4jClient neo4j.Neo4jClient) {
	clients = &clients_type{
		oktaClient,
		neo4jClient,
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
