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
	"errors"
	"fmt"
	"time"

	"github.com/intelux/insteon/plm"
	"github.com/spf13/cobra"
)

// getInfoCmd represents the on command
var getInfoCmd = &cobra.Command{
	Use:   "get-info <identity>",
	Short: "Get information about a device",
	Long:  `Get information about a device`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing identity parameter")
		}

		if len(args) > 1 {
			return errors.New("too many arguments")
		}

		identity, err := plm.ParseIdentity(args[0])

		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, time.Second)
		deviceInfo, err := powerLineModem.GetDeviceInfo(ctx, identity)

		if err != nil {
			return err
		}

		fmt.Printf("Device: %s\n", identity)
		fmt.Printf("X10 house code: %02x\n", deviceInfo.X10HouseCode)
		fmt.Printf("X10 unit: %02x\n", deviceInfo.X10Unit)
		fmt.Printf("Ramp rate: %v\n", deviceInfo.RampRate)
		fmt.Printf("On level: %.2f\n", deviceInfo.OnLevel)
		fmt.Printf("LED brightness: %.2f\n", deviceInfo.LEDBrightness)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(getInfoCmd)
}
