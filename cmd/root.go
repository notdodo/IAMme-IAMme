package cmd

import (
	"github.com/notdodo/IAMme-IAMme/pkg/infra/neo4j"
	"github.com/notdodo/IAMme-IAMme/pkg/io/logging"

	"github.com/spf13/cobra"
)

const (
	flagVerbose         = "verbose"
	flagDebug           = "debug"
	flagOrgURL          = "org-url"
	flagOktaClientToken = "client-token"
)

var (
	logger          logging.LogManager
	orgUrl          string
	oktaClientToken string
	neo4jClient     neo4j.Neo4jClient
	rootCmd         = &cobra.Command{
		Use:   "iamme",
		Short: "A CLI tool to interact with Okta and Neo4j",
	}
)

func init() {
	logger = logging.GetLogManager()
	rootCmd.PersistentFlags().BoolP(flagVerbose, "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolP(flagDebug, "d", false, "Debug output")
	rootCmd.PersistentFlags().StringVarP(&orgUrl, flagOrgURL, "u", "", "Okta Org URL")
	rootCmd.PersistentFlags().StringVarP(&oktaClientToken, flagOktaClientToken, "c", "", "Okta Client Token")
}
func Execute(neo4j neo4j.Neo4jClient) {
	neo4jClient = neo4j
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Error executing command", "err", err)
	}
}

func markAsRequired(flag string) {
	if err := rootCmd.MarkFlagRequired(flag); err != nil {
		logger.Error("Required flags not provided", "err", err, "flag", flag)
	}
}
