package ngap

// NGSetupRequestIEs Value Present
const (
	NGSetupRequestIEsValuePresentNothing int = iota
	NGSetupRequestIEsValuePresentGlobalRANNodeID
	NGSetupRequestIEsValuePresentRANNodeName
	NGSetupRequestIEsValuePresentSupportedTAList
	NGSetupRequestIEsValuePresentPagingDRX
)

// NGSetupRequestIEs Type
type NGSetupRequestIEs struct {
	ID          ProtocolIEID
	Criticality Criticality
	Value       NGSetupRequestIEsValue `aper:"openType,referenceFieldName:ID"`
}

// NGSetupRequestIEsValue OpenType
type NGSetupRequestIEsValue struct {
	Present         int
	GlobalRANNodeID *GlobalRANNodeID `aper:"referenceFieldValue:27,valueUB:3"`
	RANNodeName     *RANNodeName     `aper:"referenceFieldValue:82"`
	SupportedTAList *SupportedTAList `aper:"referenceFieldValue:102"`
	PagingDRX       *PagingDRX       `aper:"referenceFieldValue:21"`
}

// GlobalRANNodeIDExtIEs Type
type GlobalRANNodeIDExtIEs struct {
	ID          ProtocolIEID
	Criticality Criticality
	Value       GlobalRANNodeIDExtIEsValue `aper:"openType,referenceFieldName:ID"`
}

// GlobalRANNodeIDExtIEsValue OpenType
type GlobalRANNodeIDExtIEsValue struct {
	Present int
}

// GNBIDExtIEs Type
type GNBIDExtIEs struct {
	ID          ProtocolIEID
	Criticality Criticality
	Value       GNBIDExtIEsValue `aper:"openType,referenceFieldName:ID"`
}

// GNBIDExtIEsValue OpenType
type GNBIDExtIEsValue struct {
	Present int
}

// NgENBIDExtIEs Type
type NgENBIDExtIEs struct {
	ID          ProtocolIEID
	Criticality Criticality
	Value       NgENBIDExtIEsValue `aper:"openType,referenceFieldName:ID"`
}

// NgENBIDExtIEsValue OpenType
type NgENBIDExtIEsValue struct {
	Present int
}

// N3IWFIDExtIEs Type
type N3IWFIDExtIEs struct {
	ID          ProtocolIEID
	Criticality Criticality
	Value       N3IWFIDExtIEsValue `aper:"openType,referenceFieldName:ID"`
}

// N3IWFIDExtIEsValue OpenType
type N3IWFIDExtIEsValue struct {
	Present int
}
