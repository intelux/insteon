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
	"time"

	"github.com/spf13/cobra"
)

// setRampRateCmd represents the on command
var setRampRateCmd = &cobra.Command{
	Use:   "set-ramp-rate <identity> <ramp-rate>",
	Short: "Set the ramp-rate of a device",
	Long:  `Set the ramp-rate of a device`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing identity parameter")
		}

		if len(args) < 2 {
			return errors.New("missing ramp-rate parameter")
		}

		if len(args) > 2 {
			return errors.New("too many arguments")
		}

		identity, err := powerLineModem.Aliases().ParseIdentity(args[0])

		if err != nil {
			return err
		}

		rampRate, err := time.ParseDuration(args[1])

		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, time.Second)
		err = powerLineModem.SetDeviceRampRate(ctx, identity, rampRate)

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(setRampRateCmd)
}
