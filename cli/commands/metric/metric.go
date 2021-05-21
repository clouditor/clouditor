/*
 * Copyright 2021 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package metric

import (
	"clouditor.io/clouditor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewListMetricCommand returns a cobra command for the `list` subcommand
func NewListMetricCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			var client = clouditor.NewClient(viper.GetString("url"))

			err := client.StartDiscovery()

			return err
		},
	}

	return cmd
}

// NewMetricCommand returns a cobra command for `metric` subcommands
func NewMetricCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metric",
		Short: "Metric commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewListMetricCommand(),
	)
}
