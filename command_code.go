package insteon

// CommandCode represents a command code sent between the PLM and the host.
type CommandCode byte

// These types represents a command codes, as defined in the Insteon Modem
// Developer's Guide (page 12).

const (
	// Messages sent from host to PLM.
	cmdGetIMInfo                     CommandCode = 0x60
	cmdSendAllLink                   CommandCode = 0x61
	cmdSendStandardOrExtendedMessage CommandCode = 0x62
	cmdSendX10                       CommandCode = 0x63
	cmdStartAllLinking               CommandCode = 0x64
	cmdCancelAllLinking              CommandCode = 0x65
	cmdSetHostDeviceCategory         CommandCode = 0x66
	cmdResetIM                       CommandCode = 0x67
	cmdSetAckMessageByte             CommandCode = 0x68
	cmdGetFirstAllLinkRecord         CommandCode = 0x69
	cmdGetNextAllLinkRecord          CommandCode = 0x6a
	cmdSetIMConfiguration            CommandCode = 0x6b
	cmdGetAllLinkRecordForSender     CommandCode = 0x6c
	cmdLedOn                         CommandCode = 0x6d
	cmdLedOff                        CommandCode = 0x6e
	cmdManageAllLinkRecord           CommandCode = 0x6f
	cmdSetNakMessageByte             CommandCode = 0x70
	cmdSetNakMessageTwoBytes         CommandCode = 0x71
	cmdRFSleep                       CommandCode = 0x72
	cmdGetIMConfiguration            CommandCode = 0x73

	// Messages sent from PLM to host.
	cmdStandardMessageReceived     CommandCode = 0x50
	cmdExtendedMessageReceived     CommandCode = 0x51
	cmdX10Received                 CommandCode = 0x52
	cmdAllLinkingCompleted         CommandCode = 0x53
	cmdButtonEventReport           CommandCode = 0x54
	cmdUserResetDetected           CommandCode = 0x55
	cmdAllLinkCleanupFailureReport CommandCode = 0x56
	cmdAllLinkRecordMessage        CommandCode = 0x57
	cmdAllLinkCleanupStatusReport  CommandCode = 0x58
)

func isIncomingCommandCode(commandCode CommandCode) bool {
	return commandCode < cmdGetIMInfo
}

func isOutgoingCommandCode(commandCode CommandCode) bool {
	return commandCode >= cmdGetIMInfo
}
