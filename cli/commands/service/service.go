package service

import (
	"clouditor.io/clouditor/v2/cli/commands/service/assessment"
	"clouditor.io/clouditor/v2/cli/commands/service/discovery"
	"clouditor.io/clouditor/v2/cli/commands/service/orchestrator"
	"github.com/spf13/cobra"
)

// NewServiceCommand returns a cobra command for `discovery` subcommands
func NewServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Service commands",
	}

	AddCommands(cmd)

	return cmd
}

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		discovery.NewDiscoveryCommand(),
		assessment.NewAssessmentCommand(),
		orchestrator.NewOrchestratorCommand(),
	)
}
