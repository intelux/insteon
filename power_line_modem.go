package insteon

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// PowerLineModem represnts a powerline modem.
type PowerLineModem interface {
	GetIMInfo(ctx context.Context) (imInfo *IMInfo, err error)
	SetLightState(ctx context.Context, identity ID, state LightState) (err error)
	Beep(ctx context.Context, identity ID) (err error)
	GetDeviceInfo(ctx context.Context, identity ID) (deviceInfo *DeviceInfo, err error)
	SetDeviceX10Address(ctx context.Context, identity ID, x10HouseCode byte, x10Unit byte) (err error)
	SetDeviceRampRate(ctx context.Context, identity ID, rampRate time.Duration) (err error)
	SetDeviceOnLevel(ctx context.Context, identity ID, level float64) (err error)
	SetDeviceLEDBrightness(ctx context.Context, identity ID, level float64) (err error)
	GetDeviceStatus(ctx context.Context, identity ID) (level float64, err error)
	GetAllLinkDB(ctx context.Context) (records AllLinkRecordSlice, err error)
}

// DefaultPowerLineModem is the default PowerLine Modem instance.
var DefaultPowerLineModem = func() PowerLineModem {
	plm, err := NewPowerLineModem(PowerLineModemDevice)

	if err != nil {
		panic(fmt.Errorf("instanciating default PowerLine Modem: %s", err))
	}

	return plm
}()

// NewPowerLineModem instantiates a new PowerLine Modem.
func NewPowerLineModem(device string) (*SerialPowerLineModem, error) {
	url, err := url.Parse(device)

	if err != nil {
		return nil, fmt.Errorf("parsing device: %s", err)
	}

	switch url.Scheme {
	case "tcp":
		return NewRemotePowerLineModem(url.Host)
	default:
		return NewLocalPowerLineModem(url.String())
	}
}
