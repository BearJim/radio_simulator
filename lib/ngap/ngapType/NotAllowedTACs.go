package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct NotAllowedTACs */
/* TAC */
type NotAllowedTACs struct {
	List []TAC `aper:"sizeLB:1,sizeUB:16"`
}