package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct AreaOfInterestRANNodeList */
/* AreaOfInterestRANNodeItem */
type AreaOfInterestRANNodeList struct {
	List []AreaOfInterestRANNodeItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}