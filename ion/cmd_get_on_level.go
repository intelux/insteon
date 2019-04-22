package main

import (
	"fmt"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var getOnLevelCmd = &cobra.Command{
	Use:   "get-on-level <device>",
	Short: "Get the On level of a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := rootConfig.LookupDevice(args[0])

		if err != nil {
			return err
		}

		deviceInfo, err := insteon.DefaultPowerLineModem.GetDeviceInfo(rootCtx, id)

		if err != nil {
			return err
		}

		fmt.Printf("%.2f\n", deviceInfo.OnLevel)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getOnLevelCmd)
}
