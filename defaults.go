package insteon

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// EnvPowerLineModemDevice is the environment variable that contains
	// the system PowerLine Modem device.
	EnvPowerLineModemDevice = "INSTEON_POWERLINE_MODEM_DEVICE"
	// EnvPowerLineModemDebug is the environment variable that contains the
	// system PowerLine Modem debug setting.
	EnvPowerLineModemDebug = "INSTEON_POWERLINE_MODEM_DEBUG"
)

const (
	// DefaultPowerLineModemDevice is the default PowerLine Modem
	// device.
	DefaultPowerLineModemDevice = "/dev/tty0"
)

var (
	// PowerLineModemDevice is the current PowerLine Modem device.
	PowerLineModemDevice = getEnvOrDefault(
		EnvPowerLineModemDevice,
		DefaultPowerLineModemDevice,
	)

	// PowerLineModemDebug is the current PowerLine Modem debug setting.
	PowerLineModemDebug = getBoolEnvOrDefault(EnvPowerLineModemDebug, false)
)

func getEnvOrDefault(env string, def string) string {
	if value := os.Getenv(env); value != "" {
		return value
	}

	return def
}

func getBoolEnvOrDefault(env string, def bool) bool {
	if value := os.Getenv(env); value != "" {
		value = strings.ToLower(value)

		switch value {
		case "false", "no", "disabled":
			return false
		case "true", "yes", "enabled":
			return true
		default:
			v, err := strconv.Atoi(value)

			if err == nil {
				return v != 0
			}

			fmt.Fprintf(os.Stderr, "Invalid value `%s` for boolean environment variable `%s`.", value, env)
		}
	}

	return def
}
