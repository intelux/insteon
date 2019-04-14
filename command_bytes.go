package insteon

var (
	commandBytesBeep          = [2]byte{0x30, 0x00}
	commandBytesGetDeviceInfo = [2]byte{0x2e, 0x00}
	commandBytesStatusRequest = [2]byte{0x19, 0x00}
	commandBytesSetDeviceInfo = [2]byte{0x2e, 0x00}
)
