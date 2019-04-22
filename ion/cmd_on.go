package main

import (
	"errors"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var (
	onCmdInstant bool
	onCmdStep    bool
	onCmdLevel   float64
)

var onCmd = &cobra.Command{
	Use:   "on <device>",
	Short: "Turn on a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := rootConfig.LookupDevice(args[0])

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
			Level:  onCmdLevel,
			OnOff:  insteon.LightOn,
			Change: change,
		}

		if err := insteon.DefaultPowerLineModem.SetDeviceState(rootCtx, device.ID, state); err != nil {
			return err
		}

		for _, id := range device.SlaveDeviceIDs {
			insteon.DefaultPowerLineModem.SetDeviceState(rootCtx, id, state)
		}

		return nil
	},
}

func init() {
	onCmd.Flags().BoolVarP(&onCmdInstant, "instant", "i", false, "Change the light state instantly and at full value (level is ignored). Incompatible with --step.")
	onCmd.Flags().BoolVarP(&onCmdStep, "step", "s", false, "Change the light state by step (level is ignored). Incompatible with --instant.")
	onCmd.Flags().Float64VarP(&onCmdLevel, "level", "l", 1.0, "The light level, as a decimal value in the [0, 1] range.")

	rootCmd.AddCommand(onCmd)
}
