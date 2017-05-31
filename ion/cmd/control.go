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
	termbox "github.com/nsf/termbox-go"
	"github.com/spf13/cobra"
)

// controlCmd represents the control command
var controlCmd = &cobra.Command{
	Use:   "control <identity>",
	Short: "Control a device interactively",
	Long:  `Control a device interactively`,
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

		err = termbox.Init()

		if err != nil {
			return err
		}

		fmt.Println("Starting interactive control. Type <esc> to quit.")
		fmt.Println("Available controls:")
		fmt.Println("- Use arrow up/down to brighten/dim.")
		fmt.Println("- Use `O` to turn on. Ctrl+O for instant on.")
		fmt.Println("- Use `F` to turn off. Ctrl+F for instant off.")

		done := false

		setLightState := func(state plm.LightState) error {
			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, time.Second)

			return powerLineModem.SetLightState(ctx, identity, state)
		}

		changing := false

		for !done {
			event := termbox.PollEvent()

			if event.Type == termbox.EventKey {
				err = nil

				if event.Ch != 0 {
					switch event.Ch {
					case 'o':
						err = setLightState(plm.LightState{OnOff: plm.LightOn, Level: 1.0})
					case 'f':
						err = setLightState(plm.LightState{OnOff: plm.LightOff})
					}
				} else {
					switch event.Key {
					case termbox.KeyEsc:
						done = true
					case termbox.KeyArrowUp:
						if changing {
							err = setLightState(plm.LightState{OnOff: plm.LightOn, Change: plm.ChangeStop})
							changing = false
						} else {
							err = setLightState(plm.LightState{OnOff: plm.LightOn, Change: plm.ChangeStart})
							changing = true
						}
					case termbox.KeyArrowDown:
						if changing {
							err = setLightState(plm.LightState{OnOff: plm.LightOff, Change: plm.ChangeStop})
							changing = false
						} else {
							err = setLightState(plm.LightState{OnOff: plm.LightOff, Change: plm.ChangeStart})
							changing = true
						}
					case termbox.KeyCtrlO:
						err = setLightState(plm.LightState{OnOff: plm.LightOn, Change: plm.ChangeInstant})
					case termbox.KeyCtrlF:
						err = setLightState(plm.LightState{OnOff: plm.LightOff, Change: plm.ChangeInstant})
					}
				}
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(controlCmd)
}
