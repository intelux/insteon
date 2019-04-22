package main

import (
	"time"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var setRampRateCmd = &cobra.Command{
	Use:   "set-ramp-rate <identity> <level>",
	Short: "Set the ramp-rate of a device",
	Long:  `Set the ramp-rate of a device`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := rootConfig.LookupDevice(args[0])

		if err != nil {
			return err
		}

		rampRate, err := time.ParseDuration(args[1])

		if err != nil {
			return err
		}

		if err := insteon.DefaultPowerLineModem.SetDeviceRampRate(rootCtx, device.ID, rampRate); err != nil {
			return err
		}

		for _, id := range device.SlaveDeviceIDs {
			insteon.DefaultPowerLineModem.SetDeviceRampRate(rootCtx, id, rampRate)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setRampRateCmd)
}
