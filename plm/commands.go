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
	// SendAllLink sends a all-link command.
	SendAllLink CommandCode = 0x61
	// SendStandardOrExtendedMessage send a message to the PLM.
	SendStandardOrExtendedMessage = 0x62
	// SendX10 sends a X10 message.
	SendX10 = 0x63
	// StartAllLinking starts the all-linking process.
	StartAllLinking = 0x64
	// CancelAllLinking cancels the all-linking process.
	CancelAllLinking = 0x65
	// SetHostDeviceCategory sets the host device category.
	SetHostDeviceCategory = 0x66
	// ResetIM resets the PLM.
	ResetIM = 0x67
	// SetAckMessageByte sets the ack message byte.
	SetAckMessageByte = 0x68
	// GetFirstAllLinkRecord asks for the first all-link record.
	GetFirstAllLinkRecord = 0x69
	// GetNextAllLinkRecord asks for the next all-link record.
	GetNextAllLinkRecord = 0x6a
	// SetIMConfiguration sets the IM configuration.
	SetIMConfiguration = 0x6b
	// GetAllLinkRecordForSender gets the all link record for the sender.
	GetAllLinkRecordForSender = 0x6c
	// LedOn sets the IM led on.
	LedOn = 0x6d
	// LedOff sets the IM led off.
	LedOff = 0x6e
	// ManageAllLinkRecord managers the all-link records.
	ManageAllLinkRecord = 0x6f
	// SetNakMessageByte sets the Nak message byte.
	SetNakMessageByte = 0x70
	// SetNakMessageTwoBytes sets the two Nak message bytes.
	SetNakMessageTwoBytes = 0x71
	// RFSleep sets the RF antenna to sleep.
	RFSleep = 0x72
	// GetIMConfiguration gets the IM configuration.
	GetIMConfiguration = 0x73

	// Messages sent from PLM to host.

	// StandardMessageReceived is used when the PLM transmits a received
	// standard message to the host.
	StandardMessageReceived = 0x50
	// ExtendedMessageReceived is used when the PLM transmits a received
	// extended message to the host.
	ExtendedMessageReceived = 0x51
	// X10Received is used when a X10 message was received.
	X10Received = 0x52
	// AllLinkingCompleted is used when the all-linking process is completed.
	AllLinkingCompleted = 0x53
	// ButtonEventReport is used to send button event reports.
	ButtonEventReport = 0x54
	// UserResetDetected is used when a user reset is detected.
	UserResetDetected = 0x55
	// AllLinkCleanupFailureReport is used when a failure to cleanup the all-linking process is reported.
	AllLinkCleanupFailureReport = 0x56
	// AllLinkRecord is used when a all-linking record response is received.
	AllLinkRecord = 0x57
	// AllLinkCleanupStatusReport is used when a all-linking status report is received.
	AllLinkCleanupStatusReport = 0x58
)
