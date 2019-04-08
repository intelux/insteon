package plm

import (
	"context"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
)

// homekit implements a monitor that syncs with Homekit.
type homekit struct {
	Config      hc.Config
	Accessories []interface{}
	transport   hc.Transport
	mapping     map[Identity]interface{}
}

// Homekit represents a Homekit provider.
type Homekit interface {
	Initialize(plm *PowerLineModem) error
	Finalize(plm *PowerLineModem) error
}

// NewHomekit instantiates a new HomekitMonitor.
func NewHomekit(config hc.Config, accessories []interface{}) Homekit {
	return &homekit{
		Config:      config,
		Accessories: accessories,
	}
}

func (m *homekit) Initialize(plm *PowerLineModem) (err error) {
	ctx := context.Background()
	var accessories []*accessory.Accessory
	m.mapping = map[Identity]interface{}{}

	for _, acc := range m.Accessories {
		switch acc := acc.(type) {
		case *accessory.Lightbulb:
			var identity Identity

			if identity, err = plm.Aliases().ParseIdentity(acc.Info.SerialNumber.Value.(string)); err != nil {
				return
			}

			acc.Lightbulb.Brightness.OnValueRemoteUpdate(m.MakeBrightnessChangeCallback(plm, identity, acc.Lightbulb.Brightness))
			acc.Lightbulb.On.OnValueRemoteUpdate(m.MakeOnChangeCallback(plm, identity, acc.Lightbulb.Brightness))

			accessories = append(accessories, acc.Accessory)
			m.mapping[identity] = acc
		}
	}

	info := accessory.Info{
		Name:         "Ion",
		Manufacturer: "Intelux",
	}
	mainAccessory := accessory.New(info, accessory.TypeBridge)

	m.transport, err = hc.NewIPTransport(m.Config, mainAccessory, accessories...)

	if err != nil {
		return err
	}

	go m.transport.Start()

	for identity, acc := range m.mapping {
		switch acc := acc.(type) {
		case *accessory.Lightbulb:
			time.Sleep(time.Millisecond * 200)
			lvl, err := plm.GetDeviceStatus(ctx, identity)

			if err != nil {
				return err
			}

			acc.Lightbulb.Brightness.SetValue(levelToBrightness(lvl, acc.Lightbulb.Brightness))
			acc.Lightbulb.On.SetValue(lvl > 0)
		}
	}

	return nil
}

func (m *homekit) Finalize(*PowerLineModem) error {
	<-m.transport.Stop()

	return nil
}

func brightnessToLevel(brightness *characteristic.Brightness) float64 {
	return float64(brightness.GetValue()-brightness.GetMinValue()) / float64(brightness.GetMaxValue()-brightness.GetMinValue())
}

func levelToBrightness(lvl float64, brightness *characteristic.Brightness) int {
	return int(float64(brightness.GetMaxValue()-brightness.GetMinValue())*lvl) + brightness.GetMinValue()
}

func (m *homekit) LightStateUpdated(plm *PowerLineModem, id Identity, state LightState) {
	if acc := m.mapping[id]; acc != nil {
		switch acc := acc.(type) {
		case *accessory.Lightbulb:
			switch state.OnOff {
			case LightOn:
				acc.Lightbulb.On.SetValue(true)
				acc.Lightbulb.Brightness.SetValue(levelToBrightness(state.Level, acc.Lightbulb.Brightness))
			case LightOff:
				acc.Lightbulb.On.SetValue(false)
				acc.Lightbulb.Brightness.SetValue(levelToBrightness(0, acc.Lightbulb.Brightness))
			}
		}
	}
}

func (m *homekit) MakeOnChangeCallback(plm *PowerLineModem, identity Identity, brightness *characteristic.Brightness) func(bool) {
	return func(on bool) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		var state LightState

		if on {
			state = LightState{
				OnOff:  LightOn,
				Level:  brightnessToLevel(brightness),
				Change: ChangeNormal,
			}
		} else {
			state = LightState{
				OnOff:  LightOff,
				Level:  brightnessToLevel(brightness),
				Change: ChangeNormal,
			}
		}
		plm.SetLightState(ctx, identity, state)
	}
}

func (m *homekit) MakeBrightnessChangeCallback(plm *PowerLineModem, identity Identity, brightness *characteristic.Brightness) func(int) {
	return func(lvl int) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		state := LightState{
			OnOff:  LightOn,
			Level:  float64(lvl-brightness.GetMinValue()) / float64(brightness.GetMaxValue()),
			Change: ChangeNormal,
		}
		plm.SetLightState(ctx, identity, state)
	}
}
