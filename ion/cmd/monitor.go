// Copyright © 2017 Julien Kauffmann
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

package cmd

import (
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// monitorCmd represents the on command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor all activity on the device network",
	Long:  `Monitor all activity on the device network`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("monitor", true)

		return RootCmd.PersistentPreRunE(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt)

		<-stop

		return nil
	},
}

func init() {
	RootCmd.AddCommand(monitorCmd)
}
