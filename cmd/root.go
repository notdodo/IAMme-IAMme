package cmd

import (
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"

	"github.com/spf13/cobra"
)

var logger logging.LogManager
var verbose bool
var orgUrl string
var oktaClientToken string
var neo4jClient neo4j.Neo4jClient
var rootCmd = &cobra.Command{
	Use:   "iamme",
	Short: "A CLI tool to interact with Okta and Neo4j",
}

func init() {
	logger = logging.NewLogManager()
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&orgUrl, "org-url", "u", "", "Okta Org URL")
	rootCmd.PersistentFlags().StringVarP(&oktaClientToken, "client-token", "c", "", "Okta Client Token")
}

func Execute(neo4j neo4j.Neo4jClient) {
	neo4jClient = neo4j
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Error executing command", "err", err)
	}
}
