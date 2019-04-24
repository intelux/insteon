package main

import (
	"fmt"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor the PLM network activity",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		events := make(chan insteon.DeviceEvent, 10)

		go func() {
			for event := range events {
				fmt.Println(event)
			}
		}()

		insteon.DefaultPowerLineModem.Monitor(rootCtx, events)
		close(events)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
