package ngap

// NGAPPDUPresent CHOICE value
const (
	NGAPPDUPresentNothing int = iota
	NGAPPDUPresentIntiatingMessage
	NGAPPDUPresentSuccessfulOutcome
	NGAPPDUPresentUnSuccessfulOutcome
)

// NGAPPDU TYPE
type NGAPPDU struct {
	Present             int
	IntiatingMessage    *IntiatingMessage
	SuccessfulOutcome   *SuccessfulOutcome
	UnSuccessfulOutcome *UnSuccessfulOutcome
}

// TODO:use parser to generate structs
type SuccessfulOutcome struct {
}
type UnSuccessfulOutcome struct {
}
