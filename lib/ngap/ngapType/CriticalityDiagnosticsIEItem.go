package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type CriticalityDiagnosticsIEItem struct {
	IECriticality Criticality
	IEID          ProtocolIEID
	TypeOfError   TypeOfError
	IEExtensions  *ProtocolExtensionContainerCriticalityDiagnosticsIEItemExtIEs `aper:"optional"`
}
