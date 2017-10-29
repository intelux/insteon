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

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/intelux/insteon/plm"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	home             string
	cfgFile          string
	powerLineModem   *plm.PowerLineModem
	homekitTransport hc.Transport
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ion",
	Short: "Control Insteon devices",
	Long: `ion is a command-line utility that controls and monitors Insteon
devices through a local or remote PowerLine Modem device.

ion needs to be familiarized with your devices, which you can do with the help
of the "ion init" command.

Type "ion -h" to discover all the other available commands.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		device, err := plm.ParseDevice(viper.GetString("device"))

		if err != nil {
			return err
		}

		powerLineModem = plm.New(device)

		for alias, identityString := range config.Aliases {
			identity, err := plm.ParseIdentity(identityString)

			if err != nil {
				return fmt.Errorf("invalid alias value for `%s`: %s", alias, err)
			}

			powerLineModem.Aliases().Add(alias, identity)
		}

		if viper.GetBool("debug") {
			powerLineModem.SetDebugStream(os.Stderr)
		}

		var responses chan plm.Response

		if viper.GetBool("monitor") || viper.GetBool("homekit") {
			responses = make(chan plm.Response)

			go func() {
				for response := range responses {
					if viper.GetBool("monitor") {
						fmt.Println(response)
					}

					if viper.GetBool("homekit") {
						// TODO: Read the real events here.
						fmt.Println("Homekit: ", response)
					}
				}
			}()
		}

		powerLineModem.Start(responses)

		if viper.GetBool("homekit") {
			var accessories []*accessory.Accessory

			for deviceAlias, deviceType := range config.Homekit {
				info := accessory.Info{
					Name:         deviceAlias,
					SerialNumber: config.Aliases[deviceAlias],
				}

				if deviceType == "light" {
					accessory := accessory.NewLightbulb(info)
					accessory.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
						ctx := context.Background()
						identity, _ := plm.ParseIdentity(info.SerialNumber)
						var state plm.LightState

						if on {
							state = plm.LightState{
								OnOff:  plm.LightOn,
								Level:  1,
								Change: plm.ChangeNormal,
							}
						} else {
							state = plm.LightState{
								OnOff:  plm.LightOff,
								Level:  0,
								Change: plm.ChangeNormal,
							}
						}
						powerLineModem.SetLightState(ctx, identity, state)
					})
					accessories = append(accessories, accessory.Accessory)
				}
			}

			if len(accessories) == 0 {
				return fmt.Errorf("no devices are registered as Homekit devices in the configuration")
			}

			info := accessory.Info{
				Name: "Ion",
			}
			mainAccessory := accessory.New(info, accessory.TypeBridge)

			config := hc.Config{
				Pin:         "12341234",
				StoragePath: filepath.Join(home, ".config", "ion", "homekit"),
			}

			if err = os.MkdirAll(config.StoragePath, 0755); err != nil {
				return err
			}

			homekitTransport, err = hc.NewIPTransport(config, mainAccessory, accessories...)

			if err != nil {
				return err
			}

			fmt.Printf("Starting Homekit emulation for %d device(s). PIN code is: %s\n", len(accessories), config.Pin)

			go homekitTransport.Start()
		}

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if homekitTransport != nil {
			fmt.Printf("Stopping Homekit emulation.\n")

			<-homekitTransport.Stop()
		}

		if powerLineModem != nil {
			powerLineModem.Close()
		}
	},
	SilenceUsage: true,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	var err error
	home, err = homedir.Dir()

	if err != nil {
		panic(err)
	}

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ion.yaml)")
	RootCmd.PersistentFlags().String("device", "/dev/ttyUSB0", "The device to use that is connected to the PLM. Can be either a serial port or a TCP URL")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug output. For instance, this displays the RAW bytes as sent and received to/from the PLM.")
	RootCmd.PersistentFlags().BoolP("monitor", "m", false, "Enable monitoring mode. Any received response will be decoded and printed to the standard output. This differs from `debug` which does not decode the responses.")
	RootCmd.PersistentFlags().Bool("homekit", false, "Enable Homekit emulation.")

	viper.SetEnvPrefix("ion")
	viper.BindEnv("device")
	viper.BindPFlag("device", RootCmd.PersistentFlags().Lookup("device"))
	viper.BindEnv("debug")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindEnv("monitor")
	viper.BindPFlag("monitor", RootCmd.PersistentFlags().Lookup("monitor"))
	viper.BindEnv("homekit")
	viper.BindPFlag("homekit", RootCmd.PersistentFlags().Lookup("homekit"))
}

// Config describes the configuration file.
type Config struct {
	Aliases map[string]string
	Homekit map[string]string
}

var config Config

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// Search config in home directory with name ".ion" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ion")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())

		if err = viper.Unmarshal(&config); err != nil {
			panic(err)
		}
	}
}
