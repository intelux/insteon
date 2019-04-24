package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var getIMInfoCmd = &cobra.Command{
	Use:   "get-im-info",
	Short: "Get information about the PowerLine Modem itself",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		imInfo, err := insteon.DefaultPowerLineModem.GetIMInfo(rootCtx)

		if err != nil {
			return err
		}

		w := &tabwriter.Writer{}
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "ID\tCategory\tFirmware version\n")
		fmt.Fprintf(w, "%s\t%s\t%d\n", imInfo.ID, imInfo.Category, imInfo.FirmwareVersion)

		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(getIMInfoCmd)
}
