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

package main

import (
	"strings"

	"clouditor.io/clouditor/v2/cli/commands"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/clouditor/")
	viper.AddConfigPath("$HOME/.clouditor")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()
}

func newRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cl",
		Short: "The Clouditor CLI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// check, if server was specified either using config file or
			// flags, otherwise we cannot continue
			/*if viper.GetString(URLFlag) == "" {
				 return fmt.Errorf("please specify an url using the config file or the --%s flag", URLFlag)
			 }*/

			return nil
		},
	}

	commands.AddCommands(cmd)

	return cmd
}

func main() {
	var cmd = newRootCommand()

	// ignore the error returned from execute, since it will already be printed out
	_ = cmd.Execute()
}
