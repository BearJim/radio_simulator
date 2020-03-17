package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct TNLMappingList */
/* TNLMappingItem */
type TNLMappingList struct {
	List []TNLMappingItem `aper:"valueExt,sizeLB:1,sizeUB:4"`
}
