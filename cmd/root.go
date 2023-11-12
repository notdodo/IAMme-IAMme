package cmd

import (
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"
	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"

	"github.com/spf13/cobra"
)

type clientsType struct {
	okta  okta.OktaClient
	neo4j neo4j.Neo4jClient
}

var clients *clientsType
var logger logging.LogManager
var rootCmd = &cobra.Command{
	Use:   "iamme",
	Short: "A CLI tool to interact with Okta and Neo4j",
}

func init() {
	logger = logging.NewLogManager()
}

func Execute(oktaClient okta.OktaClient, neo4jClient neo4j.Neo4jClient) {
	clients = &clientsType{
		oktaClient,
		neo4jClient,
	}
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Error executing command", "err", err)
	}
}
