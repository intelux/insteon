package insteon

var (
	generalizedControllers  MainCategory = 0x00
	dimmableLightingControl MainCategory = 0x01
	switchedLightingControl MainCategory = 0x02
	networkBridges          MainCategory = 0x03
	irrigationControl       MainCategory = 0x04
	climateControlHeating   MainCategory = 0x05
	poolAndSpaControl       MainCategory = 0x06
	sensorsAndActuators     MainCategory = 0x07
	homeEntertainement      MainCategory = 0x08
	energyManagement        MainCategory = 0x09
	builtInApplianceControl MainCategory = 0x0A
	plumbing                MainCategory = 0x0B
	communication           MainCategory = 0x0C
	computerControl         MainCategory = 0x0D
	windowCoverings         MainCategory = 0x0E
	accessControl           MainCategory = 0x0F
	securityHealthSafety    MainCategory = 0x10
	surveillance            MainCategory = 0x11
	automotive              MainCategory = 0x12
	petCare                 MainCategory = 0x13
	toys                    MainCategory = 0x14
	timekeeping             MainCategory = 0x15
	holiday                 MainCategory = 0x16
	unassigned              MainCategory = 0xFF

	// networkBridges subcategories.
	powerlincSerial               SubCategory = 0x01
	powerlincUsb                  SubCategory = 0x02
	iconPowerlincSerial           SubCategory = 0x03
	iconPowerlincUsb              SubCategory = 0x04
	smartlabsPowerLineModemSerial SubCategory = 0x05
	powerlincDualBandSerial       SubCategory = 0x11
	powerlincDualBandUsb          SubCategory = 0x15
)

// MainCategory represents a main category.
type MainCategory uint8

// SubCategory represents a main category.
type SubCategory uint8

// Category represents a category.
type Category struct {
	mainCategory MainCategory
	subCategory  SubCategory
}

func (c Category) String() string {
	switch c.mainCategory {
	case generalizedControllers:
		return "Generalized Controllers"
	case dimmableLightingControl:
		return "Dimmable Lighting Control"
	case switchedLightingControl:
		return "Switched Lighting Control"
	case networkBridges:
		switch c.subCategory {
		case powerlincSerial:
			return "PowerLinc Serial [2414S]"
		case powerlincUsb:
			return "PowerLinc USB [2414U]"
		case iconPowerlincSerial:
			return "Icon PowerLinc Serial [2814 S]"
		case iconPowerlincUsb:
			return "Icon PowerLinc USB [2814U] "
		case smartlabsPowerLineModemSerial:
			return "Smartlabs Power Line Modem Serial [2412S]"
		case powerlincDualBandSerial:
			return "PowerLinc Dual Band Serial [2413S]"
		case powerlincDualBandUsb:
			return "PowerLinc Dual Band USB [2413U]"
		}

		return "Network Bridges"
	case irrigationControl:
		return "Irrigation Control"
	case climateControlHeating:
		return "Climate Control"
	case poolAndSpaControl:
		return "Pool and Spa Control"
	case sensorsAndActuators:
		return "Sensors and Actuators"
	case homeEntertainement:
		return "Home Entertainment"
	case energyManagement:
		return "Energy Management"
	case builtInApplianceControl:
		return "Built-In Appliance Control"
	case plumbing:
		return "Plumbing"
	case communication:
		return "Communication"
	case computerControl:
		return "Computer Control"
	case windowCoverings:
		return "Window Coverings"
	case accessControl:
		return "Access Control"
	case securityHealthSafety:
		return "Security Health Safety"
	case surveillance:
		return "Surveillance"
	case automotive:
		return "Automotive"
	case petCare:
		return "Pet Care"
	case toys:
		return "Toys"
	case timekeeping:
		return "Timekeeping"
	case holiday:
		return "Holiday"
	case unassigned:
		return "Unassigned"
	}

	return "Unknown category"
}

// MarshalText marshals a category as text.
func (c Category) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}
