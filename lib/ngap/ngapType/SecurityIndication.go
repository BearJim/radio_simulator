package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type SecurityIndication struct {
	IntegrityProtectionIndication       IntegrityProtectionIndication
	ConfidentialityProtectionIndication ConfidentialityProtectionIndication
	MaximumIntegrityProtectedDataRate   *MaximumIntegrityProtectedDataRate                  `aper:"optional"`
	IEExtensions                        *ProtocolExtensionContainerSecurityIndicationExtIEs `aper:"optional"`
}
