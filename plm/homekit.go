package plm

import (
	"context"
	"fmt"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
)

// homekitMonitor implements a monitor that syncs with Homekit.
type homekitMonitor struct {
	Config      hc.Config
	Accessories []interface{}
	transport   hc.Transport
}

// NewHomekitMonitor instantiates a new HomekitMonitor.
func NewHomekitMonitor(config hc.Config, accessories []interface{}) (monitor Monitor) {
	monitor = &homekitMonitor{
		Config:      config,
		Accessories: accessories,
	}

	return
}

func (m *homekitMonitor) Initialize(plm *PowerLineModem) (err error) {
	var accessories []*accessory.Accessory

	for _, acc := range m.Accessories {
		switch acc := acc.(type) {
		case *accessory.Lightbulb:
			var identity Identity

			if identity, err = plm.Aliases().ParseIdentity(acc.Info.Name.Value.(string)); err != nil {
				return
			}

			acc.Lightbulb.Brightness.OnValueRemoteUpdate(m.MakeBrightnessChangeCallback(plm, identity, acc.Lightbulb.Brightness))
			acc.Lightbulb.On.OnValueRemoteUpdate(m.MakeOnChangeCallback(plm, identity, acc.Lightbulb.Brightness))
			accessories = append(accessories, acc.Accessory)
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

	return nil
}

func (m *homekitMonitor) Finalize(*PowerLineModem) error {
	<-m.transport.Stop()

	return nil
}

func (m *homekitMonitor) ResponseReceived(plm *PowerLineModem, response Response) {
	switch response := response.(type) {
	case *StandardMessageReceivedResponse:
		if lightState := CommandBytesToLightState(response.CommandBytes); lightState != nil {
			m.LightStateUpdated(plm, response.Target, *lightState)
		}
	}
}

func (m *homekitMonitor) LightStateUpdated(plm *PowerLineModem, id Identity, state LightState) {
	// TODO: Call something like: acc.Lightbulb.On.SetValue(true)
	fmt.Println(id, state)
}

func (m *homekitMonitor) MakeOnChangeCallback(plm *PowerLineModem, identity Identity, brightness *characteristic.Brightness) func(bool) {
	return func(on bool) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		var state LightState

		if on {
			state = LightState{
				OnOff:  LightOn,
				Level:  float64(brightness.GetValue()-brightness.GetMinValue()) / float64(brightness.GetMaxValue()),
				Change: ChangeNormal,
			}
		} else {
			state = LightState{
				OnOff:  LightOff,
				Level:  float64(brightness.GetValue()-brightness.GetMinValue()) / float64(brightness.GetMaxValue()),
				Change: ChangeNormal,
			}
		}
		plm.SetLightState(ctx, identity, state)
	}
}

func (m *homekitMonitor) MakeBrightnessChangeCallback(plm *PowerLineModem, identity Identity, brightness *characteristic.Brightness) func(int) {
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
