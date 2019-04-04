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
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information about the PLM",
	Long:  `Displays information about the PowerLine Modem device.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, time.Second)
		info, err := powerLineModem.GetInfo(ctx)

		if err != nil {
			return err
		}

		fmt.Printf("Category: %s\n", info.Category)
		fmt.Printf("Identity: %s\n", info.Identity)
		fmt.Printf("Firmware version: %d\n", info.FirmwareVersion)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
