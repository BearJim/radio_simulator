package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct QosFlowAddOrModifyRequestList */
/* QosFlowAddOrModifyRequestItem */
type QosFlowAddOrModifyRequestList struct {
	List []QosFlowAddOrModifyRequestItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}