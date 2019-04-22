package insteon

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// WebService implements a web-service that manages Insteon devices.
type WebService struct {
	PowerLineModem PowerLineModem
	Configuration  *Configuration

	once    sync.Once
	handler http.Handler
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

func (s *WebService) init() {
	s.once.Do(func() {
		if s.PowerLineModem == nil {
			s.PowerLineModem = DefaultPowerLineModem
		}

		if s.Configuration == nil {
			s.Configuration = &Configuration{}
		}

		s.handler = s.makeHandler()
	})
}

func (s *WebService) makeHandler() http.Handler {
	router := mux.NewRouter()

	router.Path("/api/im-info").Methods(http.MethodGet).HandlerFunc(s.handleGetIMInfo)
	router.Path("/api/all-link-db").Methods(http.MethodGet).HandlerFunc(s.handleGetAllLinkDB)
	router.Path("/api/device/{id}/state").Methods(http.MethodGet).HandlerFunc(s.handleGetDeviceState)
	router.Path("/api/device/{id}/state").Methods(http.MethodPut).HandlerFunc(s.handleSetDeviceState)
	router.Path("/api/device/{id}/info").Methods(http.MethodGet).HandlerFunc(s.handleGetDeviceInfo)
	router.Path("/api/device/{id}/info").Methods(http.MethodPut).HandlerFunc(s.handleSetDeviceInfo)
	router.Path("/api/device/{id}/beep").Methods(http.MethodPost).HandlerFunc(s.handleBeep)

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
	device := s.parseDevice(w, r)

	if device == nil {
		return
	}

	state, err := s.PowerLineModem.GetDeviceState(r.Context(), device.ID)

	if err != nil {
		s.handleError(w, r, err)
		return
	}

	s.handleValue(w, r, state)
}

func (s *WebService) handleSetDeviceState(w http.ResponseWriter, r *http.Request) {
	device := s.parseDevice(w, r)

	if device == nil {
		return
	}

	state := &LightState{}

	if !s.decodeValue(w, r, state) {
		return
	}

	if err := s.PowerLineModem.SetDeviceState(r.Context(), device.ID, *state); err != nil {
		s.handleError(w, r, err)
		return
	}

	s.handleValue(w, r, state)
}

func (s *WebService) handleBeep(w http.ResponseWriter, r *http.Request) {
	device := s.parseDevice(w, r)

	if device == nil {
		return
	}

	if err := s.PowerLineModem.Beep(r.Context(), device.ID); err != nil {
		s.handleError(w, r, err)
		return
	}
}

func (s *WebService) handleGetDeviceInfo(w http.ResponseWriter, r *http.Request) {
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

func (s *WebService) handleSetDeviceInfo(w http.ResponseWriter, r *http.Request) {
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

func (s *WebService) parseDevice(w http.ResponseWriter, r *http.Request) *ConfigurationDevice {
	vars := mux.Vars(r)

	id := vars["id"]

	if id == "" {
		err := fmt.Errorf("invalid empty device id")

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
