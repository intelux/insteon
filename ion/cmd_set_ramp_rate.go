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
		id, err := insteon.ParseID(args[0])

		if err != nil {
			return err
		}

		rampRate, err := time.ParseDuration(args[1])

		if err != nil {
			return err
		}

		return insteon.DefaultPowerLineModem.SetDeviceRampRate(rootCtx, id, rampRate)
	},
}

func init() {
	rootCmd.AddCommand(setRampRateCmd)
}
