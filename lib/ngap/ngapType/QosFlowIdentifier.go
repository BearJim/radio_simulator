package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type QosFlowIdentifier struct {
	Value int64 `aper:"valueExt,valueLB:0,valueUB:63"`
}
