package plm

import (
	"context"
	"fmt"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

// homekitMonitor implements a monitor that syncs with Homekit.
type homekitMonitor struct {
	Context        context.Context
	PowerLineModem *PowerLineModem
	Config         hc.Config
	Accessories    []*accessory.Accessory
	transport      hc.Transport
}

// NewHomekitMonitor instantiates a new HomekitMonitor.
func NewHomekitMonitor(ctx context.Context, plm *PowerLineModem, config hc.Config, accessories []interface{}) Monitor {
	var rawAccessories []*accessory.Accessory

	for _, acc := range accessories {
		switch acc := acc.(type) {
		case *accessory.Lightbulb:
			acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(lvl int) {
				if identity, err := plm.Aliases().ParseIdentity(acc.Info.Name.Value.(string)); err == nil {
					fmt.Println(acc.Info.Name.Value.(string), identity, err)
					plm.SetDeviceOnLevel(ctx, identity, float64(lvl-acc.Lightbulb.Brightness.GetMinValue())/float64(acc.Lightbulb.Brightness.GetMaxValue()))
				}
			})
			acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
				fmt.Println(acc.Lightbulb.Brightness.GetMinValue(), acc.Lightbulb.Brightness.GetValue(), acc.Lightbulb.Brightness.GetMaxValue())
			})
			rawAccessories = append(rawAccessories, acc.Accessory)
		}
	}

	return &homekitMonitor{
		Context:        ctx,
		PowerLineModem: plm,
		Config:         config,
		Accessories:    rawAccessories,
	}
}

func (m *homekitMonitor) Initialize() (err error) {
	info := accessory.Info{
		Name:         "Ion",
		Manufacturer: "Intelux",
	}
	mainAccessory := accessory.New(info, accessory.TypeBridge)

	m.transport, err = hc.NewIPTransport(m.Config, mainAccessory, m.Accessories...)

	if err != nil {
		return err
	}

	go m.transport.Start()

	return nil
}

func (m *homekitMonitor) Finalize() error {
	<-m.transport.Stop()

	return nil
}

func (m *homekitMonitor) LightStateUpdated(id Identity, state LightState) {
	fmt.Println(id, state)
}
