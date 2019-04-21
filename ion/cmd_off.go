package main

import (
	"errors"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var (
	offCmdInstant bool
	offCmdStep    bool
)

var offCmd = &cobra.Command{
	Use:   "off <device>",
	Short: "Turn off a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := rootConfig.LookupDevice(args[0])

		if err != nil {
			return err
		}

		var change = insteon.ChangeNormal

		if onCmdInstant {
			if onCmdStep {
				return errors.New("can't specify both `--instant` and `--step`")
			}

			change = insteon.ChangeInstant
		} else if onCmdStep {
			change = insteon.ChangeStep
		}

		state := insteon.LightState{
			Level:  0,
			OnOff:  insteon.LightOff,
			Change: change,
		}

		return insteon.DefaultPowerLineModem.SetLightState(rootCtx, id, state)
	},
}

func init() {
	offCmd.Flags().BoolVarP(&offCmdInstant, "instant", "i", false, "Change the light state instantly. Incompatible with --step.")
	offCmd.Flags().BoolVarP(&offCmdStep, "step", "s", false, "Change the light state by step. Incompatible with --instant.")

	rootCmd.AddCommand(offCmd)
}
