package ngap

import (
	"radio_simulator/lib/aper"
)

// PLMNIdentity Type
type PLMNIdentity struct {
	Value aper.OctetString `aper:"sizeLB:3,sizeUB:3"`
}
