package main

import (
	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var offCmd = &cobra.Command{
	Use:  "off <device>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := insteon.ParseID(args[0])

		if err != nil {
			return err
		}

		state := insteon.LightState{
			Level: 0,
			OnOff: insteon.LightOff,
		}

		return insteon.DefaultPowerLineModem.SetLightState(rootCtx, id, state)
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}
