package main

import (
	"encoding/json"
	"os"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var getDeviceInfoCmd = &cobra.Command{
	Use:  "get-device-info <device>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := insteon.ParseID(args[0])

		if err != nil {
			return err
		}

		deviceInfo, err := insteon.DefaultPowerLineModem.GetDeviceInfo(rootCtx, id)

		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(deviceInfo)
	},
}

func init() {
	rootCmd.AddCommand(getDeviceInfoCmd)
}
