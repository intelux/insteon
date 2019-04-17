package main

import (
	"strconv"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var setOnLevelCmd = &cobra.Command{
	Use:   "set-on-level <identity> <level>",
	Short: "Set the On level of a device",
	Long:  `Set the On level of a device`,
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

		return insteon.DefaultPowerLineModem.SetDeviceOnLevel(rootCtx, id, level)
	},
}

func init() {
	rootCmd.AddCommand(setOnLevelCmd)
}
