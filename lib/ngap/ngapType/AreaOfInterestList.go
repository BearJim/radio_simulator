package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct AreaOfInterestList */
/* AreaOfInterestItem */
type AreaOfInterestList struct {
	List []AreaOfInterestItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}
