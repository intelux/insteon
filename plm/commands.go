package plm

const (
	MessageStart   byte = 0x02
	MessageFailure byte = 0x15
)

// These types represents a command codes, as defined in the Insteon Modem
// Developer's Guide (page 12).

// RequestCode is a message sent from the host to the modem.
type RequestCode byte

// ResponseCode is a message sent from the modem to the host.
type ResponseCode byte

const (
	// GetIMInfo asks the modem for its information.
	GetIMInfo RequestCode = 0x60
)
