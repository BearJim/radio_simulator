package ngap

// BroadcastPLMNList Type
type BroadcastPLMNList struct {
	List []BroadcastPLMNItem `aper:"valueExt,sizeLB:1,sizeUB:12"`
}
