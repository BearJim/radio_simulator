package ngapType

import "radio_simulator/lib/aper"

// Need to import "radio_simulator/lib/aper" if it uses "aper"

const (
	ReflectiveQosAttributePresentSubjectTo aper.Enumerated = 0
)

type ReflectiveQosAttribute struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
