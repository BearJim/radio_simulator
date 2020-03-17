package ngap

// SupportedTAList Type
type SupportedTAList struct {
	List []SupportedTAItem `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
