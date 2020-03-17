package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type SecurityContext struct {
	NextHopChainingCount NextHopChainingCount
	NextHopNH            SecurityKey
	IEExtensions         *ProtocolExtensionContainerSecurityContextExtIEs `aper:"optional"`
}
