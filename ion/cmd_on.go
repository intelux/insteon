package main

import (
	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var onCmd = &cobra.Command{
	Use:  "on <device>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := insteon.ParseID(args[0])

		if err != nil {
			return err
		}

		state := insteon.LightState{
			Level: 1,
			OnOff: insteon.LightOn,
		}

		return insteon.DefaultPowerLineModem.SetLightState(rootCtx, id, state)
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
}
