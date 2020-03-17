package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type GTPTunnel struct {
	TransportLayerAddress TransportLayerAddress
	GTPTEID               GTPTEID
	IEExtensions          *ProtocolExtensionContainerGTPTunnelExtIEs `aper:"optional"`
}
