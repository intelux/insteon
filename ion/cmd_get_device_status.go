package main

import (
	"encoding/json"
	"os"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var getDeviceStatusCmd = &cobra.Command{
	Use:  "get-device-status <device>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := insteon.ParseID(args[0])

		if err != nil {
			return err
		}

		level, err := insteon.DefaultPowerLineModem.GetDeviceStatus(rootCtx, id)

		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(level)
	},
}

func init() {
	rootCmd.AddCommand(getDeviceStatusCmd)
}
