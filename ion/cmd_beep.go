package main

import (
	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var beepCmd = &cobra.Command{
	Use:   "beep <device>",
	Short: "Make a device beep",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := rootConfig.LookupDevice(args[0])

		if err != nil {
			return err
		}

		if err := insteon.DefaultPowerLineModem.Beep(rootCtx, device.ID); err != nil {
			return err
		}

		for _, id := range device.SlaveDeviceIDs {
			insteon.DefaultPowerLineModem.Beep(rootCtx, id)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(beepCmd)
}
