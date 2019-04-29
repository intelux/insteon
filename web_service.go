package insteon

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// WebService implements a web-service that manages Insteon devices.
type WebService struct {
	PowerLineModem PowerLineModem
	Configuration  *Configuration

	// DisablePowerLineModem is a boolean value that, if set, disables
	// exposition of the PowerLineModem routes.
	DisablePowerLineModem bool

	// DisableAPI is a boolean value that, if set, disables exposition of the
	// API routes.
	DisableAPI bool

	// ForceRefreshPeriod is the period after which to consider a device state
	// stale.
	ForceRefreshPeriod time.Duration

	once                   sync.Once
	lock                   sync.Mutex
	handler                http.Handler
	controllers            map[ID]bool
	responders             map[ID]bool
	deviceToMasterDevice   map[ID]ID
	deviceStates           map[ID]*LightState
	deviceStatesTimestamps map[ID]time.Time
}

// NewWebService instanciates a new web service.
//
// If no PowerLine modem is specified, the default one is taken.
// If no configuration is specified, the default one is taken.
func NewWebService(powerLineModem PowerLineModem, configuration *Configuration) *WebService {
	return &WebService{
		PowerLineModem: powerLineModem,
		Configuration:  configuration,
	}
}

// Handler returns the HTTP handler associated to the web-service.
func (s *WebService) Handler() http.Handler {
	s.init()

	return s.handler
}

// Synchronize the web-service with the Insteon network to optimize efficiency.
//
// Must be done before Run is called or the HTTP handler is served.
func (s *WebService) Synchronize(ctx context.Context, failOnMissingOptimization bool) error {
	s.init()

	// Read the All-Link DB to make sure we only deal with devices that we can
	// control/respond to.
	records, err := s.PowerLineModem.GetAllLinkDB(ctx)

	if err != nil {
		return err
	}

	s.controllers = map[ID]bool{}
	s.responders = map[ID]bool{}

	for _, record := range records {
		// The device is not known, we skip it.
		if _, ok := s.deviceToMasterDevice[record.ID]; !ok {
			continue
		}

		if record.Mode() == ModeResponder {
			s.responders[record.ID] = true
		} else {
			s.controllers[record.ID] = true
		}
	}

	if failOnMissingOptimization {
		var failures []string

		for _, device := range s.Configuration.Devices {
			if !s.responders[device.ID] {
				failures = append(failures, fmt.Sprintf("device %s (%s) is not a responder", device.Name, device.ID))
			}

			if !s.controllers[device.ID] {
				failures = append(failures, fmt.Sprintf("device %s (%s) is not a controller", device.Name, device.ID))
			}

			for _, mirrorDeviceID := range device.MirrorDeviceIDs {
				if !s.responders[mirrorDeviceID] {
					failures = append(failures, fmt.Sprintf("mirror device %s of %s (%s) is not a responder", mirrorDeviceID, device.Name, device.ID))
				}

				if !s.controllers[mirrorDeviceID] {
					failures = append(failures, fmt.Sprintf("mirror device %s of %s (%s) is not a controller", mirrorDeviceID, device.Name, device.ID))
				}
			}
		}

		if len(failures) > 0 {
			return fmt.Errorf("optimization failure:\n- %s", strings.Join(failures, "\n- "))
		}
	}

	return nil
}

// Run the web-service for as long as the specified context remains valid.
func (s *WebService) Run(ctx context.Context) error {
	s.init()

	events := make(chan DeviceEvent, 10)
	defer close(events)

	go func() {
		for event := range events {
			id := event.Identity

			if masterID, ok := s.deviceToMasterDevice[id]; ok {
				id = masterID
			}

			// If an event occurs for a device, remove it's cached state.
			s.lock.Lock()
			delete(s.deviceStates, id)
			delete(s.deviceStatesTimestamps, id)
			s.lock.Unlock()
		}
	}()

	return s.PowerLineModem.Monitor(ctx, events)
}

func (s *WebService) init() {
	s.once.Do(func() {
		if s.PowerLineModem == nil {
			s.PowerLineModem = DefaultPowerLineModem
		}

		if s.Configuration == nil {
			s.Configuration = &Configuration{}
		}

		if s.ForceRefreshPeriod == 0 {
			s.ForceRefreshPeriod = 5 * time.Minute
		}

		s.handler = s.makeHandler()
		s.deviceToMasterDevice = map[ID]ID{}
		s.deviceStates = map[ID]*LightState{}
		s.deviceStatesTimestamps = map[ID]time.Time{}

		for _, device := range s.Configuration.Devices {
			s.deviceToMasterDevice[device.ID] = device.ID

			for _, mirrorDeviceID := range device.MirrorDeviceIDs {
				s.deviceToMasterDevice[mirrorDeviceID] = device.ID
			}

			for _, controllerID := range device.ControllerIDs {
				s.deviceToMasterDevice[controllerID] = device.ID
			}
		}
	})
}

func (s *WebService) makeHandler() http.Handler {
	router := mux.NewRouter()

	// PLM-specific routes.
	if !s.DisablePowerLineModem {
		router.Path("/plm/im-info").Methods(http.MethodGet).HandlerFunc(s.handleGetIMInfo)
		router.Path("/plm/all-link-db").Methods(http.MethodGet).HandlerFunc(s.handleGetAllLinkDB)
		router.Path("/plm/device/{id}/state").Methods(http.MethodGet).HandlerFunc(s.handleGetDeviceState)
		router.Path("/plm/device/{id}/state").Methods(http.MethodPut).HandlerFunc(s.handleSetDeviceState)
		router.Path("/plm/device/{id}/info").Methods(http.MethodGet).HandlerFunc(s.handleGetDeviceInfo)
		router.Path("/plm/device/{id}/info").Methods(http.MethodPut).HandlerFunc(s.handleSetDeviceInfo)
		router.Path("/plm/device/{id}/beep").Methods(http.MethodPost).HandlerFunc(s.handleBeep)
	}

	// API routes.
	if !s.DisableAPI {
		router.Path("/api/device/{device}/state").Methods(http.MethodGet).HandlerFunc(s.handleAPIGetDeviceState)
		router.Path("/api/device/{device}/state").Methods(http.MethodPut).HandlerFunc(s.handleAPISetDeviceState)
		router.Path("/api/device/{device}/info").Methods(http.MethodGet).HandlerFunc(s.handleAPIGetDeviceInfo)
		router.Path("/api/device/{device}/info").Methods(http.MethodPut).HandlerFunc(s.handleAPISetDeviceInfo)
	}

	return router
}

func (s *WebService) handleGetIMInfo(w http.ResponseWriter, r *http.Request) {
	imInfo, err := s.PowerLineModem.GetIMInfo(r.Context())

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.handleValue(w, r, imInfo)
}

func (s *WebService) handleGetDeviceState(w http.ResponseWriter, r *http.Request) {
	id := s.parseID(w, r)

	if id == nil {
		return
	}

	state, err := s.PowerLineModem.GetDeviceState(r.Context(), *id)

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.handleValue(w, r, state)
}

func (s *WebService) handleSetDeviceState(w http.ResponseWriter, r *http.Request) {
	id := s.parseID(w, r)

	if id == nil {
		return
	}

	state := &LightState{}

	if !s.decodeValue(w, r, state) {
		return
	}

	if err := s.PowerLineModem.SetDeviceState(r.Context(), *id, *state); err != nil {
		s.handleError(w, r, err)
		return
	}

	s.handleValue(w, r, state)
}

func (s *WebService) handleBeep(w http.ResponseWriter, r *http.Request) {
	id := s.parseID(w, r)

	if id == nil {
		return
	}

	if err := s.PowerLineModem.Beep(r.Context(), *id); err != nil {
		s.handleError(w, r, err)
		return
	}
}

func (s *WebService) handleGetDeviceInfo(w http.ResponseWriter, r *http.Request) {
	id := s.parseID(w, r)

	if id == nil {
		return
	}

	deviceInfo, err := s.PowerLineModem.GetDeviceInfo(r.Context(), *id)

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.handleValue(w, r, deviceInfo)
}

func (s *WebService) handleSetDeviceInfo(w http.ResponseWriter, r *http.Request) {
	id := s.parseID(w, r)

	if id == nil {
		return
	}

	deviceInfo := &DeviceInfo{}

	if !s.decodeValue(w, r, deviceInfo) {
		return
	}

	if err := s.PowerLineModem.SetDeviceInfo(r.Context(), *id, *deviceInfo); err != nil {
		s.handleError(w, r, err)
		return
	}

	s.handleValue(w, r, deviceInfo)
}

func (s *WebService) handleGetAllLinkDB(w http.ResponseWriter, r *http.Request) {
	records, err := s.PowerLineModem.GetAllLinkDB(r.Context())

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.handleValue(w, r, records)
}

func (s *WebService) handleAPIGetDeviceState(w http.ResponseWriter, r *http.Request) {
	device := s.parseDevice(w, r)

	if device == nil {
		return
	}

	now := time.Now().UTC()

	s.lock.Lock()
	state := s.deviceStates[device.ID]

	if state != nil {
		timestamp := s.deviceStatesTimestamps[device.ID]

		if timestamp.Add(s.ForceRefreshPeriod).Before(now) {
			delete(s.deviceStatesTimestamps, device.ID)
			delete(s.deviceStates, device.ID)
			state = nil
		}
	}
	s.lock.Unlock()

	if state == nil {
		var err error
		state, err = s.PowerLineModem.GetDeviceState(r.Context(), device.ID)

		if err != nil {
			s.handleError(w, r, err)
			return
		}

		// Only cache the state if the device is a controller, otherwise it
		// won't ever be refreshed.
		if s.controllers != nil && s.controllers[device.ID] {
			s.lock.Lock()
			s.deviceStates[device.ID] = state
			s.deviceStatesTimestamps[device.ID] = time.Now().UTC()
			s.lock.Unlock()
		}
	}

	s.handleValue(w, r, state)
}

func (s *WebService) handleAPISetDeviceState(w http.ResponseWriter, r *http.Request) {
	device := s.parseDevice(w, r)

	if device == nil {
		return
	}

	state := &LightState{}

	if !s.decodeValue(w, r, state) {
		return
	}

	// If the device is not a responder, don't bother sending a command to it.
	if s.responders != nil && !s.responders[device.ID] {
		err := fmt.Errorf("device %s (%s) is registered as a responder", device.Name, device.ID)
		s.handleError(w, r, err)
		return
	}

	if err := s.PowerLineModem.SetDeviceState(r.Context(), device.ID, *state); err != nil {
		s.handleError(w, r, err)
		return
	}

	for _, id := range device.MirrorDeviceIDs {
		if s.responders == nil || s.responders[id] {
			s.PowerLineModem.SetDeviceState(r.Context(), id, *state)
		}
	}

	// Only cache the state if the device is a controller, otherwise it
	// won't ever be refreshed.
	if s.controllers != nil && s.controllers[device.ID] {
		s.lock.Lock()
		s.deviceStates[device.ID] = state
		s.deviceStatesTimestamps[device.ID] = time.Now().UTC()
		s.lock.Unlock()
	}

	s.handleValue(w, r, state)
}

func (s *WebService) handleAPIGetDeviceInfo(w http.ResponseWriter, r *http.Request) {
	device := s.parseDevice(w, r)

	if device == nil {
		return
	}

	deviceInfo, err := s.PowerLineModem.GetDeviceInfo(r.Context(), device.ID)

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.handleValue(w, r, deviceInfo)
}

func (s *WebService) handleAPISetDeviceInfo(w http.ResponseWriter, r *http.Request) {
	device := s.parseDevice(w, r)

	if device == nil {
		return
	}

	deviceInfo := &DeviceInfo{}

	if !s.decodeValue(w, r, deviceInfo) {
		return
	}

	if err := s.PowerLineModem.SetDeviceInfo(r.Context(), device.ID, *deviceInfo); err != nil {
		s.handleError(w, r, err)
		return
	}

	for _, id := range device.MirrorDeviceIDs {
		s.PowerLineModem.SetDeviceInfo(r.Context(), id, *deviceInfo)
	}

	s.handleValue(w, r, deviceInfo)
}

func (s *WebService) parseID(w http.ResponseWriter, r *http.Request) *ID {
	vars := mux.Vars(r)

	idStr := vars["id"]

	if idStr == "" {
		err := fmt.Errorf("invalid empty device id")

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)

		return nil
	}

	id, err := ParseID(idStr)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)

		return nil
	}

	return &id
}

func (s *WebService) parseDevice(w http.ResponseWriter, r *http.Request) *ConfigurationDevice {
	vars := mux.Vars(r)

	id := vars["device"]

	if id == "" {
		err := fmt.Errorf("invalid empty device identifier")

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)

		return nil
	}

	device, err := s.Configuration.LookupDevice(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)

		return nil
	}

	return device
}

func (s *WebService) decodeValue(w http.ResponseWriter, r *http.Request, value interface{}) bool {
	if r.Body == nil {
		err := fmt.Errorf("missing body")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)

		return false
	}

	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)

		return false
	}

	switch mediatype {
	case "":
		mediatype = "application/json"
	case "application/json":
	default:
		err := fmt.Errorf("expected body of type application/json")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)

		return false
	}

	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)

		return false
	}

	return true
}

func (s *WebService) handleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", err)
}

func (s *WebService) handleValue(w http.ResponseWriter, r *http.Request, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(value)
}
