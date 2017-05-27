package plm

const (
	// MessageStart is the marker at the beginning of commands.
	MessageStart byte = 0x02
	// MessageAck is returned as an acknowledgment.
	MessageAck byte = 0x06
	// MessageNak is returned as an non-acknowledgment.
	MessageNak byte = 0x15
)

// These types represents a command codes, as defined in the Insteon Modem
// Developer's Guide (page 12).

// CommandCode represents a command code sent between the PLM and the host.
type CommandCode byte

const (
	// Messages sent from host to PLM.

	// GetIMInfo asks the modem for its information.
	GetIMInfo CommandCode = 0x60
	// SendStandardOrExtendedMessage send a message to the PLM.
	SendStandardOrExtendedMessage = 0x62

	// Messages sent from PLM to host.

	// StandardMessageReceived is used when the PLM transmits a received
	// standard message to the host.
	StandardMessageReceived = 0x50
	// ExtendedMessageReceived is used when the PLM transmits a received
	// extended message to the host.
	ExtendedMessageReceived = 0x51
)
