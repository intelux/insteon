package main

import (
	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var beepCmd = &cobra.Command{
	Use:  "beep <device>",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := insteon.ParseID(args[0])

		if err != nil {
			return err
		}

		return insteon.DefaultPowerLineModem.Beep(rootCtx, id)
	},
}

func init() {
	rootCmd.AddCommand(beepCmd)
}
