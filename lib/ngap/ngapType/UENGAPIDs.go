package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	UENGAPIDsPresentNothing int = iota /* No components present */
	UENGAPIDsPresentUENGAPIDPair
	UENGAPIDsPresentAMFUENGAPID
	UENGAPIDsPresentChoiceExtensions
)

type UENGAPIDs struct {
	Present          int
	UENGAPIDPair     *UENGAPIDPair `aper:"valueExt"`
	AMFUENGAPID      *AMFUENGAPID
	ChoiceExtensions *ProtocolIESingleContainerUENGAPIDsExtIEs
}
