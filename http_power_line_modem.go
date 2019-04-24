package insteon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"sync"
)

// HTTPPowerLineModem implements a PowerLine modem over HTTP.
type HTTPPowerLineModem struct {
	URL    *url.URL
	Client *http.Client

	once sync.Once
}

// NewHTTPPowerLineModem instanciates a new HTTP PowerLine modem.
func NewHTTPPowerLineModem(p string) (*HTTPPowerLineModem, error) {
	u, err := url.Parse(p)

	if err != nil {
		return nil, fmt.Errorf("parsing URL: %s", err)
	}

	return &HTTPPowerLineModem{
		URL: u,
	}, nil
}

// GetIMInfo gets information about the PowerLine Modem.
func (m *HTTPPowerLineModem) GetIMInfo(ctx context.Context) (imInfo *IMInfo, err error) {
	imInfo = &IMInfo{}
	err = m.do(ctx, http.MethodGet, "/plm/im-info", nil, imInfo)

	return
}

// GetAllLinkDB gets the on level of a device.
func (m *HTTPPowerLineModem) GetAllLinkDB(ctx context.Context) (records AllLinkRecordSlice, err error) {
	err = m.do(ctx, http.MethodGet, "/plm/all-link-db", nil, &records)

	return
}

// GetDeviceState gets the on level of a device.
func (m *HTTPPowerLineModem) GetDeviceState(ctx context.Context, identity ID) (state *LightState, err error) {
	url := fmt.Sprintf("/plm/device/%s/state", identity)
	state = &LightState{}
	err = m.do(ctx, http.MethodGet, url, nil, state)

	return
}

// SetDeviceState sets the state of a lighting device.
func (m *HTTPPowerLineModem) SetDeviceState(ctx context.Context, identity ID, state LightState) error {
	url := fmt.Sprintf("/plm/device/%s/state", identity)

	return m.do(ctx, http.MethodPut, url, state, nil)
}

// GetDeviceInfo returns the information about a device.
func (m *HTTPPowerLineModem) GetDeviceInfo(ctx context.Context, identity ID) (deviceInfo *DeviceInfo, err error) {
	url := fmt.Sprintf("/plm/device/%s/info", identity)
	deviceInfo = &DeviceInfo{}
	err = m.do(ctx, http.MethodGet, url, nil, deviceInfo)

	return
}

// SetDeviceInfo sets the information on device.
func (m *HTTPPowerLineModem) SetDeviceInfo(ctx context.Context, identity ID, deviceInfo DeviceInfo) error {
	url := fmt.Sprintf("/plm/device/%s/info", identity)

	return m.do(ctx, http.MethodPut, url, nil, deviceInfo)
}

// Beep causes a device to beep.
func (m *HTTPPowerLineModem) Beep(ctx context.Context, identity ID) error {
	url := fmt.Sprintf("/plm/device/%s/beep", identity)

	return m.do(ctx, http.MethodPost, url, nil, nil)
}

// Monitor the Insteon network for changes for as long as the specified context remains valid.
//
// All events are pushed to the specified events channel.
func (m *HTTPPowerLineModem) Monitor(ctx context.Context, events chan<- DeviceEvent) error {
	event := DeviceEvent{}

	for {
		if err := m.do(ctx, http.MethodGet, "/plm/device/next-event", nil, &event); err == nil {
			select {
			case events <- event:
			case <-ctx.Done():
				return ctx.Err()
			}
		} else {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}
	}
}

func (m *HTTPPowerLineModem) init() {
	m.once.Do(func() {
		if m.URL == nil {
			m.URL = &url.URL{
				Scheme: "http",
				Host:   "localhost:7660",
			}
		}

		if m.Client == nil {
			m.Client = http.DefaultClient
		}
	})
}

func (m *HTTPPowerLineModem) do(ctx context.Context, method string, path string, input interface{}, output interface{}) error {
	m.init()

	var body io.Reader

	if input != nil {
		buf := &bytes.Buffer{}

		if err := json.NewEncoder(buf).Encode(input); err != nil {
			return fmt.Errorf("encoding request body: %s", err)
		}

		body = buf
	}

	u := *m.URL
	u.Path = path

	req, err := http.NewRequest(method, u.String(), body)

	if err != nil {
		return fmt.Errorf("creating new request: %s", err)
	}

	req = req.WithContext(ctx)

	if input != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := m.Client.Do(req)

	if err != nil {
		return fmt.Errorf("executing request: %s", err)
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if output != nil {
		mediatype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))

		if err != nil {
			return fmt.Errorf("parsing response content-type: %s", err)
		}

		switch mediatype {
		case "":
			mediatype = "application/json"
		case "application/json":
		default:
			return fmt.Errorf("expected body of type application/json")
		}

		if err = json.NewDecoder(resp.Body).Decode(output); err != nil {
			return fmt.Errorf("decoding response: %s", err)
		}
	}

	return nil
}
