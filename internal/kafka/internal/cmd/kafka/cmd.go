// Package cluster contains commands for interacting with cluster logic of the service directly instead of through the
// REST API exposed via the serve command.
package kafka

import (
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/environments"
	"github.com/spf13/cobra"
)

func NewKafkaCommand(env *environments.Env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kafka",
		Short: "Perform kafka CRUD actions directly",
		Long:  "Perform kafka CRUD actions directly.",
	}

	// add sub-commands
	cmd.AddCommand(
		NewCreateCommand(env),
		NewGetCommand(env),
		NewDeleteCommand(env),
		NewListCommand(env),
	)

	return cmd
}