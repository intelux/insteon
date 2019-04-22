package main

import (
	"strconv"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var setOnLevelCmd = &cobra.Command{
	Use:   "set-on-level <identity> <level>",
	Short: "Set the On level of a device",
	Long:  `Set the On level of a device`,
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

		if err := insteon.DefaultPowerLineModem.SetDeviceOnLevel(rootCtx, device.ID, level); err != nil {
			return err
		}

		for _, id := range device.SlaveDeviceIDs {
			insteon.DefaultPowerLineModem.SetDeviceOnLevel(rootCtx, id, level)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setOnLevelCmd)
}
