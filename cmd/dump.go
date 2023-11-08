package cmd

import (
	"github.com/notdodo/IAMme-IAMme/pkg/app"

	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "dump",
	Short: "Fetch Okta info and store them in Neo4j",
	Run: func(cmd *cobra.Command, args []string) {
		oktaNeo4jApp := app.NewOktaNeo4jApp(clients.okta, clients.neo4j)
		oktaNeo4jApp.Dump()
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
