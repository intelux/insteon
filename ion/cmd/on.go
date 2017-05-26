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
	onInstant bool
	onLevel   float64
)

// onCmd represents the on command
var onCmd = &cobra.Command{
	Use:   "on <identity>",
	Short: "Turn a light on",
	Long:  `Turn a light on`,
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
		state := plm.LightState{
			OnOff:   plm.LightOn,
			Instant: onInstant,
			Level:   onLevel,
		}
		err = powerLineModem.SetLightState(ctx, identity, state)

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	onCmd.Flags().BoolVarP(&onInstant, "instant", "i", false, "Change the light state instantly.")
	onCmd.Flags().Float64VarP(&onLevel, "level", "l", 1.0, "The light level, as a decimal value.")
	RootCmd.AddCommand(onCmd)
}
