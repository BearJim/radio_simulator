package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct CriticalityDiagnostics_IE_List */
/* CriticalityDiagnosticsIEItem */
type CriticalityDiagnosticsIEList struct {
	List []CriticalityDiagnosticsIEItem `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
