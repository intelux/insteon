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
		id, err := insteon.ParseID(args[0])

		if err != nil {
			return err
		}

		deviceInfo, err := insteon.DefaultPowerLineModem.GetDeviceInfo(rootCtx, id)

		if err != nil {
			return err
		}

		fmt.Printf("%s\n", deviceInfo.RampRate)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getRampRateCmd)
}
