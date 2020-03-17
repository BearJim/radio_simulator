package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	ResetTypePresentNothing int = iota /* No components present */
	ResetTypePresentNGInterface
	ResetTypePresentPartOfNGInterface
	ResetTypePresentChoiceExtensions
)

type ResetType struct {
	Present           int
	NGInterface       *ResetAll
	PartOfNGInterface *UEAssociatedLogicalNGConnectionList
	ChoiceExtensions  *ProtocolIESingleContainerResetTypeExtIEs
}
