package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceAdmittedList */
/* PDUSessionResourceAdmittedItem */
type PDUSessionResourceAdmittedList struct {
	List []PDUSessionResourceAdmittedItem `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
