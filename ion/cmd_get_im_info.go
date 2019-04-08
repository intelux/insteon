package main

import (
	"encoding/json"
	"os"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var getIMInfoCmd = &cobra.Command{
	Use: "get-im-info",
	RunE: func(cmd *cobra.Command, args []string) error {
		imInfo, err := insteon.DefaultPowerLineModem.GetIMInfo(rootCtx)

		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(imInfo)
	},
}

func init() {
	rootCmd.AddCommand(getIMInfoCmd)
}
