package main

import (
	"strconv"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var setLevelCmd = &cobra.Command{
	Use:   "set-level <identity> <level>",
	Short: "Set the level of a device",
	Long:  `Set the level of a device`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := insteon.ParseID(args[0])

		if err != nil {
			return err
		}

		level, err := strconv.ParseFloat(args[1], 64)

		if err != nil {
			return err
		}

		state := insteon.LightState{
			Level:  level,
			OnOff:  insteon.LightOn,
			Change: insteon.ChangeInstant,
		}

		return insteon.DefaultPowerLineModem.SetLightState(rootCtx, id, state)
	},
}

func init() {
	rootCmd.AddCommand(setLevelCmd)
}
