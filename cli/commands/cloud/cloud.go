// Copyright 2021 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package cloud

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/emptypb"
)

// NewRegisterCertificationTargetCommand returns a cobra command for the `discover` subcommand
func NewRegisterCertificationTargetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [name]",
		Short: "Registers a new target certification target",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.CertificationTarget
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			id := args[0]
			name := args[1]

			res, err = client.RegisterCertificationTarget(context.Background(), &orchestrator.RegisterCertificationTargetRequest{
				CertificationTarget: &orchestrator.CertificationTarget{
					Id:   id,
					Name: name,
				},
			})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"aws", "azure"}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewListCertificationTargetsCommand returns a cobra command for the `list` subcommand
func NewListCertificationTargetsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all target certification targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.ListCertificationTargetsResponse
				target  []*orchestrator.CertificationTarget
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			target, err = api.ListAllPaginated(&orchestrator.ListCertificationTargetsRequest{}, client.ListCertificationTargets, func(res *orchestrator.ListCertificationTargetsResponse) []*orchestrator.CertificationTarget {
				return res.Targets
			})

			// Build a response with all certification targets
			res = &orchestrator.ListCertificationTargetsResponse{
				Targets: target,
			}

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewGetCertificationTargetCommand returns a cobra command for the `get` subcommand
func NewGetCertificationTargetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Retrieves a target certification target by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.CertificationTarget
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			targetID := args[0]

			res, err = client.GetCertificationTarget(context.Background(), &orchestrator.GetCertificationTargetRequest{CertificationTargetId: targetID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	return cmd
}

// NewRemoveCertificationTargetComand returns a cobra command for the `get` subcommand
func NewRemoveCertificationTargetComand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [id]",
		Short: "Removes a target certification target by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *emptypb.Empty
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			targetID := args[0]

			res, err = client.RemoveCertificationTarget(context.Background(), &orchestrator.RemoveCertificationTargetRequest{CertificationTargetId: targetID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	return cmd
}

// NewUpdateCertificationTargetCommand returns a cobra command for the `update` subcommand
func NewUpdateCertificationTargetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a target certification target",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.CertificationTarget
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			res, err = client.UpdateCertificationTarget(context.Background(), &orchestrator.UpdateCertificationTargetRequest{
				CertificationTarget: &orchestrator.CertificationTarget{
					Id:          viper.GetString("id"),
					Name:        viper.GetString("name"),
					Description: viper.GetString("description"),
				},
			})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	cmd.PersistentFlags().String("id", "", "the certification target id to update")
	cmd.PersistentFlags().StringP("name", "n", "", "the name of the certification target")
	cmd.PersistentFlags().StringP("description", "d", "", "an optional description")

	_ = cmd.MarkPersistentFlagRequired("id")
	_ = cmd.MarkPersistentFlagRequired("name")
	_ = viper.BindPFlag("id", cmd.PersistentFlags().Lookup("id"))
	_ = viper.BindPFlag("name", cmd.PersistentFlags().Lookup("name"))
	_ = viper.BindPFlag("description", cmd.PersistentFlags().Lookup("description"))
	_ = viper.BindPFlag("control-ids", cmd.PersistentFlags().Lookup("control-ids"))

	_ = cmd.RegisterFlagCompletionFunc("id", cli.ValidArgsGetCertificationTargets)
	_ = cmd.RegisterFlagCompletionFunc("name", cli.DefaultArgsShellComp)
	_ = cmd.RegisterFlagCompletionFunc("description", cli.DefaultArgsShellComp)

	return cmd
}

// NewGetMetricConfigurationCommand returns a cobra command for the `get-metric-configuration` subcommand
// TODO(oxisto): Can we have something like cl cloud get <id> metric-configuration <id>?
func NewGetMetricConfigurationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-metric-configuration",
		Short: "Retrieves a metric configuration for a specific certification target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *assessment.MetricConfiguration
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			targetID := args[0]
			metricID := args[1]

			res, err = client.GetMetricConfiguration(context.Background(), &orchestrator.GetMetricConfigurationRequest{CertificationTargetId: targetID, MetricId: metricID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewCloudCommand returns a cobra command for `cloud` subcommands
func NewCloudCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud",
		Short: "Target certification targets commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewRegisterCertificationTargetCommand(),
		NewListCertificationTargetsCommand(),
		NewGetCertificationTargetCommand(),
		NewUpdateCertificationTargetCommand(),
		NewRemoveCertificationTargetComand(),
		NewGetMetricConfigurationCommand(),
	)
}
