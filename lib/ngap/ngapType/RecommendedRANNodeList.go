package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct RecommendedRANNodeList */
/* RecommendedRANNodeItem */
type RecommendedRANNodeList struct {
	List []RecommendedRANNodeItem `aper:"valueExt,sizeLB:1,sizeUB:16"`
}
