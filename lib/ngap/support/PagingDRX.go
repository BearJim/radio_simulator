package ngap

import "radio_simulator/lib/aper"

// PagingDRX ENUMURATED
type PagingDRX struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:3"`
}

// PagingDRX ENUMURATED Value
const (
	PagingDRXv32 aper.Enumerated = iota
	PagingDRXv64
	PagingDRXv128
	PagingDRXv256
)
