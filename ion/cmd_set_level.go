package main

import (
	"strconv"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var setLevelCmd = &cobra.Command{
	Use:   "set-level <identity> <level>",
	Short: "Set the level of a device",
	Long:  `Set the level of a device`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := rootConfig.LookupDevice(args[0])

		if err != nil {
			return err
		}

		level, err := strconv.ParseFloat(args[1], 64)

		if err != nil {
			return err
		}

		state := insteon.LightState{
			Level:  level,
			OnOff:  insteon.LightOn,
			Change: insteon.ChangeNormal,
		}

		if err := insteon.DefaultPowerLineModem.SetLightState(rootCtx, device.ID, state); err != nil {
			return err
		}

		for _, id := range device.SlaveDeviceIDs {
			insteon.DefaultPowerLineModem.SetLightState(rootCtx, id, state)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setLevelCmd)
}
