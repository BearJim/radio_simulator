package ngap

// GlobalRANNodeIDPresent CHOICE value
const (
	GlobalRANNodeIDPresentNothing int = iota
	GlobalRANNodeIDPresentGlobalGNBID
	GlobalRANNodeIDPresentGlobalNgENBID
	GlobalRANNodeIDPresentGlobalN3IWFID
	GlobalRANNodeIDPresentChoiceExtensions
)

// GlobalRANNodeID CHOICE Type
type GlobalRANNodeID struct {
	Present          int
	GlobalGNBID      *GlobalGNBID   `aper:"valueExt"`
	GlobalNgENBID    *GlobalNgENBID `aper:"valueExt"`
	GlobalN3IWFID    *GlobalN3IWFID `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerGlobalRANNodeID
}
