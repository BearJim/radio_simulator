package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceListCxtRelReq */
/* PDUSessionResourceItemCxtRelReq */
type PDUSessionResourceListCxtRelReq struct {
	List []PDUSessionResourceItemCxtRelReq `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
