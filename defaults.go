package insteon

import "os"

const (
	// EnvPowerLineModemDevice is the environment variable that contains
	// the system PowerLine Modem device.
	EnvPowerLineModemDevice = "INSTEON_POWERLINE_MODEM_DEVICE"
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
)

func getEnvOrDefault(env string, def string) string {
	if value := os.Getenv(env); value != "" {
		return value
	}

	return def
}
