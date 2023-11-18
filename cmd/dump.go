package cmd

import (
	"github.com/notdodo/IAMme-IAMme/pkg/app"
	"github.com/notdodo/IAMme-IAMme/pkg/infra/okta"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "dump",
	Short: "Fetch Okta info and store them in Neo4j",
	Run: func(cmd *cobra.Command, args []string) {
		markAsRequired("org-url")
		markAsRequired("client-token")
		if err := rootCmd.ValidateRequiredFlags(); err != nil {
			logger.Error("Required flags not provided", "err", err)
		}
		if cmd.Flags().Changed(flagVerbose) {
			logger.SetVerboseLevel()
		}
		if cmd.Flags().Changed(flagDebug) {
			logger.SetDebugLevel()
		}
		oktaNeo4jApp := app.NewOktaNeo4jApp(okta.NewOktaClient(orgUrl, oktaClientToken), neo4jClient)
		oktaNeo4jApp.Dump()
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
