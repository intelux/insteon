package hub

import (
	"context"
	"sync"
	"time"

	"github.com/intelux/insteon/plm"
)

// A Hub implements a virtual device that maintains state of real physical
// devices.
type Hub struct {
	modem       *plm.PowerLineModem
	deviceInfos []*DeviceInfo
}

// NewHub instantiates a new Hub.
func NewHub(modem *plm.PowerLineModem, devices []Device) *Hub {
	deviceInfos := make([]*DeviceInfo, len(devices))

	for i, device := range devices {
		deviceInfos[i] = NewDeviceInfo(device)
	}

	return &Hub{
		modem:       modem,
		deviceInfos: deviceInfos,
	}
}

// Now returns the current time.
func (h *Hub) Now() time.Time {
	return time.Now().UTC()
}

// Run the hub update logic.
func (h *Hub) Run(ctx context.Context) {
	period := time.Millisecond * 500

	wg := sync.WaitGroup{}
	wg.Add(len(h.deviceInfos))

	for _, deviceInfo := range h.deviceInfos {
		go func(deviceInfo *DeviceInfo) {
			defer wg.Done()

			timer := time.NewTimer(period)

			for {
				select {
				case level := <-deviceInfo.nextLevel:
					state := plm.LightState{
						OnOff:  plm.LightOn,
						Change: plm.ChangeNormal,
						Level:  level,
					}

					if err := h.modem.SetLightState(ctx, deviceInfo.ID, state); err != nil {
						// Try to reset the level, but if it fails, something
						// tried to set a new level already so we give up.
						select {
						case deviceInfo.nextLevel <- level:
						default:
						}
					} else {
						deviceInfo.State = &State{
							Level:     level,
							UpdatedAt: h.Now(),
						}
					}

					if !timer.Stop() {
						<-timer.C
					}

					timer.Reset(period)

				case <-timer.C:
					level, err := h.modem.GetDeviceStatus(ctx, deviceInfo.ID)

					if err != nil {
						continue
					}

					deviceInfo.State = &State{
						Level:     level,
						UpdatedAt: h.Now(),
					}

					timer.Reset(period)

				case <-ctx.Done():
					return
				}
			}
		}(deviceInfo)
	}

	wg.Wait()
}

// SetDeviceLevel set a device level.
func (h *Hub) SetDeviceLevel(ctx context.Context, deviceName string, level float64) {
	for _, deviceInfo := range h.deviceInfos {
		if deviceInfo.Name == deviceName {
			select {
			// If there is already a level waiting to be set and we can
			// remove it, we do so and replace it with a our new level.
			case <-deviceInfo.nextLevel:
				deviceInfo.nextLevel <- level

			// If we can insert our level, our job is done.
			case deviceInfo.nextLevel <- level:

			// If the context expires, we exit immediately.
			case <-ctx.Done():
			}
		}
	}
}

// State describes a device state.
type State struct {
	Level     float64   `json:"level"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeviceInfo represents a device and its information.
type DeviceInfo struct {
	Device
	State *State `json:"physical_state"`

	nextLevel chan float64
}

// NewDeviceInfo instantiates a new device info.
func NewDeviceInfo(device Device) *DeviceInfo {
	return &DeviceInfo{
		Device:    device,
		nextLevel: make(chan float64, 1),
	}
}
