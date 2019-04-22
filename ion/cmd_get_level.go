package main

import (
	"fmt"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var getLevelCmd = &cobra.Command{
	Use:   "get-level <device>",
	Short: "Get the current level of a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := rootConfig.LookupDevice(args[0])

		if err != nil {
			return err
		}

		level, err := insteon.DefaultPowerLineModem.GetDeviceStatus(rootCtx, device.ID)

		if err != nil {
			return err
		}

		fmt.Printf("%.2f\n", level)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getLevelCmd)
}
