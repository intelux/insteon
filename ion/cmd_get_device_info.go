package main

import (
	"encoding/hex"
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
		device, err := rootConfig.LookupDevice(args[0])

		if err != nil {
			return err
		}

		deviceInfo, err := insteon.DefaultPowerLineModem.GetDeviceInfo(rootCtx, device.ID)

		if err != nil {
			return err
		}

		w := &tabwriter.Writer{}
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "Attribute\tValue\n")
		fmt.Fprintf(w, "X10 Address\t%s\n", hex.EncodeToString((*deviceInfo.X10Address)[:]))
		fmt.Fprintf(w, "Ramp rate\t%s\n", *deviceInfo.RampRate)
		fmt.Fprintf(w, "On level\t%.2f\n", *deviceInfo.OnLevel)
		fmt.Fprintf(w, "LED brightness\t%.2f\n", *deviceInfo.LEDBrightness)
		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(getDeviceInfoCmd)
}
