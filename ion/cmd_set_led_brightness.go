package main

import (
	"strconv"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var setLEDBrightnessCmd = &cobra.Command{
	Use:   "set-led-brightness <identity> <level>",
	Short: "Set the led-brightness of a device",
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

		deviceInfo := insteon.DeviceInfo{
			LEDBrightness: &level,
		}

		if err := insteon.DefaultPowerLineModem.SetDeviceInfo(rootCtx, device.ID, deviceInfo); err != nil {
			return err
		}

		for _, id := range device.MirrorDeviceIDs {
			insteon.DefaultPowerLineModem.SetDeviceInfo(rootCtx, id, deviceInfo)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setLEDBrightnessCmd)
}
