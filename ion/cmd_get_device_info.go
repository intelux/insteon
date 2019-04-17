package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var getDeviceInfoCmd = &cobra.Command{
	Use:   "get-device-info <device>",
	Short: "Get all the available information about a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := insteon.ParseID(args[0])

		if err != nil {
			return err
		}

		deviceInfo, err := insteon.DefaultPowerLineModem.GetDeviceInfo(rootCtx, id)

		if err != nil {
			return err
		}

		level, err := insteon.DefaultPowerLineModem.GetDeviceStatus(rootCtx, id)

		if err != nil {
			return err
		}

		w := &tabwriter.Writer{}
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "Attribute\tValue\n")
		fmt.Fprintf(w, "X10 House Code\t%02x\n", deviceInfo.X10HouseCode)
		fmt.Fprintf(w, "X10 Unit\t%02x\n", deviceInfo.X10Unit)
		fmt.Fprintf(w, "Ramp rate\t%s\n", deviceInfo.RampRate)
		fmt.Fprintf(w, "On level\t%.2f\n", deviceInfo.OnLevel)
		fmt.Fprintf(w, "LED brightness\t%.2f\n", deviceInfo.LEDBrightness)
		fmt.Fprintf(w, "Level\t%.2f\n", level)
		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(getDeviceInfoCmd)
}
