package ngap

// GlobalGNBIDExtIEs Type
type GlobalGNBIDExtIEs struct {
	ID             ProtocolIEID
	Criticality    Criticality
	ExtensionValue GlobalGNBIDExtIEsValue `aper:"openType,referenceFieldType:ID"`
}

// GlobalGNBIDExtIEsValue Type
type GlobalGNBIDExtIEsValue struct {
	Present int
}

// GlobalNgENBIDExtIEs Type
type GlobalNgENBIDExtIEs struct {
	ID             ProtocolIEID
	Criticality    Criticality
	ExtensionValue GlobalNgENBIDExtIEsValue `aper:"openType,referenceFieldType:ID"`
}

// GlobalNgENBIDExtIEsValue Type
type GlobalNgENBIDExtIEsValue struct {
	Present int
}

// GlobalN3IWFIDExtIEs Type
type GlobalN3IWFIDExtIEs struct {
	ID             ProtocolIEID
	Criticality    Criticality
	ExtensionValue GlobalN3IWFIDExtIEsValue `aper:"openType,referenceFieldType:ID"`
}

// GlobalN3IWFIDExtIEsValue Type
type GlobalN3IWFIDExtIEsValue struct {
	Present int
}

// SupportedTAItemExtIEs Type
type SupportedTAItemExtIEs struct {
	ID             ProtocolIEID
	Criticality    Criticality
	ExtensionValue SupportedTAItemExtIEsValue `aper:"openType,referenceFieldType:ID"`
}

// SupportedTAItemExtIEsValue Type
type SupportedTAItemExtIEsValue struct {
	Present int
}

// BroadcastPLMNItemExtIEs Type
type BroadcastPLMNItemExtIEs struct {
	ID             ProtocolIEID
	Criticality    Criticality
	ExtensionValue BroadcastPLMNItemExtIEsValue `aper:"openType,referenceFieldType:ID"`
}

// BroadcastPLMNItemExtIEsValue Type
type BroadcastPLMNItemExtIEsValue struct {
	Present int
}

// SupportItemExtIEs Type
type SupportItemExtIEs struct {
	ID             ProtocolIEID
	Criticality    Criticality
	ExtensionValue SupportItemExtIEsValue `aper:"openType,referenceFieldType:ID"`
}

// SupportItemExtIEsValue Type
type SupportItemExtIEsValue struct {
	Present int
}

// SNSSAIItemExtIEs Type
type SNSSAIItemExtIEs struct {
	ID             ProtocolIEID
	Criticality    Criticality
	ExtensionValue SNSSAIItemExtIEsValue `aper:"openType,referenceFieldType:ID"`
}

// SNSSAIItemExtIEsValue Type
type SNSSAIItemExtIEsValue struct {
	Present int
}
