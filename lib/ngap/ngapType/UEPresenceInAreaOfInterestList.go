package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct UEPresenceInAreaOfInterestList */
/* UEPresenceInAreaOfInterestItem */
type UEPresenceInAreaOfInterestList struct {
	List []UEPresenceInAreaOfInterestItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}
