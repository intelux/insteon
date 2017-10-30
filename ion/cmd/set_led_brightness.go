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
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// setLEDBrightnessCmd represents the on command
var setLEDBrightnessCmd = &cobra.Command{
	Use:   "set-led-brightness <identity> <level>",
	Short: "Set the led-brightness of a device",
	Long:  `Set the led-brightness of a device`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing identity parameter")
		}

		if len(args) < 2 {
			return errors.New("missing on-level parameter")
		}

		if len(args) > 2 {
			return errors.New("too many arguments")
		}

		identity, err := powerLineModem.Aliases().ParseIdentity(args[0])

		if err != nil {
			return err
		}

		level, err := strconv.ParseFloat(args[1], 64)

		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		err = powerLineModem.SetDeviceLEDBrightness(ctx, identity, level)

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(setLEDBrightnessCmd)
}
