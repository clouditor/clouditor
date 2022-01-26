package service

import (
	"clouditor.io/clouditor/cli/commands/service/assessment"
	"clouditor.io/clouditor/cli/commands/service/discovery"
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
	)
}
