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

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/emptypb"
)

// NewRegisterCloudServiceCommand returns a cobra command for the `discover` subcommand
func NewRegisterCloudServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [name]",
		Short: "Registers a new target cloud service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.CloudService
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			name := args[0]

			res, err = client.RegisterCloudService(context.Background(), &orchestrator.RegisterCloudServiceRequest{
				Service: &orchestrator.CloudService{
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

// NewListCloudServicesCommand returns a cobra command for the `list` subcommand
func NewListCloudServicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all target cloud services",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.ListCloudServicesResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			res, err = client.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewGetCloudServiceCommand returns a cobra command for the `get` subcommand
func NewGetCloudServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Retrieves a target cloud service by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.CloudService
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			serviceID := args[0]

			res, err = client.GetCloudService(context.Background(), &orchestrator.GetCloudServiceRequest{ServiceId: serviceID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	return cmd
}

// NewRemoveCloudServiceComand returns a cobra command for the `get` subcommand
func NewRemoveCloudServiceComand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [id]",
		Short: "Removes a target cloud service by its ID",
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

			serviceID := args[0]

			res, err = client.RemoveCloudService(context.Background(), &orchestrator.RemoveCloudServiceRequest{ServiceId: serviceID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	return cmd
}

// NewUpdateCloudServiceCommand returns a cobra command for the `update` subcommand
func NewUpdateCloudServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a target cloud service",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.CloudService
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			res, err = client.UpdateCloudService(context.Background(), &orchestrator.UpdateCloudServiceRequest{
				ServiceId: viper.GetString("id"),
				Service: &orchestrator.CloudService{
					Name:        viper.GetString("name"),
					Description: viper.GetString("description"),
				},
			})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.ValidArgsGetCloudServices,
	}

	cmd.PersistentFlags().String("id", "", "the cloud service id to update")
	cmd.PersistentFlags().StringP("name", "n", "", "the name of the service")
	cmd.PersistentFlags().StringP("description", "d", "", "an optional description")
	_ = cmd.MarkPersistentFlagRequired("id")
	_ = cmd.MarkPersistentFlagRequired("name")
	_ = viper.BindPFlag("id", cmd.PersistentFlags().Lookup("id"))
	_ = viper.BindPFlag("name", cmd.PersistentFlags().Lookup("name"))
	_ = viper.BindPFlag("description", cmd.PersistentFlags().Lookup("description"))

	_ = cmd.RegisterFlagCompletionFunc("id", cli.ValidArgsGetCloudServices)
	_ = cmd.RegisterFlagCompletionFunc("name", cli.DefaultArgsShellComp)
	_ = cmd.RegisterFlagCompletionFunc("description", cli.DefaultArgsShellComp)

	return cmd
}

// NewGetMetricConfigurationCommand returns a cobra command for the `get-metric-configuration` subcommand
// TODO(oxisto): Can we have something like cl cloud get <id> metric-configuration <id>?
func NewGetMetricConfigurationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-metric-configuration",
		Short: "Retrieves a metric configuration for a specific target service",
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

			serviceID := args[0]
			metricID := args[1]

			res, err = client.GetMetricConfiguration(context.Background(), &orchestrator.GetMetricConfigurationRequest{ServiceId: serviceID, MetricId: metricID})

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
		Short: "Target cloud services commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewRegisterCloudServiceCommand(),
		NewListCloudServicesCommand(),
		NewGetCloudServiceCommand(),
		NewUpdateCloudServiceCommand(),
		NewRemoveCloudServiceComand(),
		NewGetMetricConfigurationCommand(),
	)
}
