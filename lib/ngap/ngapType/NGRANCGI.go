package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	NGRANCGIPresentNothing int = iota /* No components present */
	NGRANCGIPresentNRCGI
	NGRANCGIPresentEUTRACGI
	NGRANCGIPresentChoiceExtensions
)

type NGRANCGI struct {
	Present          int
	NRCGI            *NRCGI    `aper:"valueExt"`
	EUTRACGI         *EUTRACGI `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerNGRANCGIExtIEs
}
