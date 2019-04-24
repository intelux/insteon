package insteon

import (
	"context"
	"fmt"
	"net/url"
)

// PowerLineModem represnts a powerline modem.
type PowerLineModem interface {
	GetIMInfo(ctx context.Context) (imInfo *IMInfo, err error)
	GetAllLinkDB(ctx context.Context) (records AllLinkRecordSlice, err error)
	GetDeviceState(ctx context.Context, identity ID) (state *LightState, err error)
	SetDeviceState(ctx context.Context, identity ID, state LightState) (err error)
	GetDeviceInfo(ctx context.Context, identity ID) (deviceInfo *DeviceInfo, err error)
	SetDeviceInfo(ctx context.Context, identity ID, deviceInfo DeviceInfo) error
	Beep(ctx context.Context, identity ID) (err error)
	Monitor(ctx context.Context, events chan<- DeviceEvent) error
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
func NewPowerLineModem(device string) (PowerLineModem, error) {
	url, err := url.Parse(device)

	if err != nil {
		return nil, fmt.Errorf("parsing device: %s", err)
	}

	switch url.Scheme {
	case "http", "https":
		return NewHTTPPowerLineModem(url.String())
	case "tcp":
		return NewRemotePowerLineModem(url.Host)
	default:
		return NewLocalPowerLineModem(url.String())
	}
}
