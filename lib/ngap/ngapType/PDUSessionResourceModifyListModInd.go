package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceModifyListModInd */
/* PDUSessionResourceModifyItemModInd */
type PDUSessionResourceModifyListModInd struct {
	List []PDUSessionResourceModifyItemModInd `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
