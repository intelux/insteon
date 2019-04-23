package main

import (
	"fmt"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var getRampRateCmd = &cobra.Command{
	Use:   "get-ramp-rate <device>",
	Short: "Get the ramp-rate of a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := rootConfig.LookupDevice(args[0])

		if err != nil {
			return err
		}

		deviceInfo, err := insteon.DefaultPowerLineModem.GetDeviceInfo(rootCtx, device.ID)

		if err != nil {
			return err
		}

		fmt.Printf("%s\n", *deviceInfo.RampRate)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getRampRateCmd)
}
