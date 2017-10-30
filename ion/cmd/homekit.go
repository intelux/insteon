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
	"os"
	"path/filepath"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/intelux/insteon/plm"
	"github.com/spf13/cobra"
)

// homekitCmd represents the on command
var homekitCmd = &cobra.Command{
	Use:   "homekit",
	Short: "Start a Homekit emulator for all the known devices.",
	Long:  "Start a Homekit emulator for all the known devices.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		var accessories []interface{}

		for deviceAlias, deviceType := range config.Homekit {
			info := accessory.Info{
				Name: deviceAlias,
			}

			if deviceType == "light" {
				accessories = append(accessories, accessory.NewLightbulb(info))
			}
		}

		homekitConfig := hc.Config{
			Pin:         "12341234",
			StoragePath: filepath.Join(home, ".config", "ion", "homekit"),
		}

		if err := os.MkdirAll(homekitConfig.StoragePath, 0755); err != nil {
			return err
		}

		homekitMonitor := plm.NewHomekitMonitor(ctx, powerLineModem, homekitConfig, accessories)

		fmt.Printf("Starting Homekit emulation for %d device(s). PIN code is: %s\n", len(config.Homekit), homekitConfig.Pin)

		return powerLineModem.Monitor(ctx, homekitMonitor)
	},
}

func init() {
	RootCmd.AddCommand(homekitCmd)
}
