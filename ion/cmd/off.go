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

	"github.com/intelux/insteon/plm"
	"github.com/spf13/cobra"
)

var (
	offInstant bool
	offStep    bool
)

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off <identity>",
	Short: "Turn a light off",
	Long:  `Turn a light off`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing identity parameter")
		}

		if len(args) > 1 {
			return errors.New("too many arguments")
		}

		identity, err := powerLineModem.Aliases().ParseIdentity(args[0])

		if err != nil {
			return err
		}

		var change = plm.ChangeNormal

		if offInstant {
			if offStep {
				return errors.New("can't specify both `--instant` and `--step`")
			}

			change = plm.ChangeInstant
		} else if offStep {
			change = plm.ChangeStep
		}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, time.Second)
		state := plm.LightState{
			OnOff:  plm.LightOff,
			Change: change,
			Level:  0,
		}
		err = powerLineModem.SetLightState(ctx, identity, state)

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	offCmd.Flags().BoolVarP(&offInstant, "instant", "i", false, "Change the light state instantly. Incompatible with --step.")
	offCmd.Flags().BoolVarP(&offStep, "step", "s", false, "Change the light state by step. Incompatible with --instant.")
	RootCmd.AddCommand(offCmd)
}
