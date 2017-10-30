package plm

// Monitor represents a type that can listen on device changes.
type Monitor interface {
	Initialize() error
	Finalize() error
	LightStateUpdated(id Identity, state LightState)
}
