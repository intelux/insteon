// Copyright © 2017 Julien Kauffmann
// {{.copyright}}
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"os"
	"time"

	"github.com/intelux/insteon/hub"
	"github.com/intelux/insteon/plm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		device, err := plm.ParseDevice(os.Getenv("ION_DEVICE"))

		if err != nil {
			return err
		}

		powerLineModem := plm.New(device)
		powerLineModem.SetDebugStream(os.Stderr)

		if err = powerLineModem.Start(nil); err != nil {
			return err
		}

		defer powerLineModem.Close()

		devices := []hub.Device{
			hub.Device{
				ID:   plm.Identity{0x40, 0x99, 0x15},
				Name: "hall",
			},
			hub.Device{
				ID:   plm.Identity{0x42, 0x7f, 0x3e},
				Name: "kitchen",
			},
			hub.Device{
				ID:   plm.Identity{0x40, 0x98, 0x73},
				Name: "stairs",
			},
		}

		hub := hub.NewHub(powerLineModem, devices)

		ctx := context.Background()

		go func() {
			time.Sleep(time.Second * 3)
			hub.SetDeviceLevel(ctx, "kitchen", 0.5)
			time.Sleep(time.Second * 2)
			hub.SetDeviceLevel(ctx, "kitchen", 1)
			hub.SetDeviceLevel(ctx, "hall", 1)
			time.Sleep(time.Second * 2)
			hub.SetDeviceLevel(ctx, "kitchen", 0)
			hub.SetDeviceLevel(ctx, "hall", 0)
			time.Sleep(time.Second * 2)
		}()

		hub.Run(ctx)

		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
