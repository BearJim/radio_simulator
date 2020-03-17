package ngap

// ProcedureCode type
type ProcedureCode struct {
	Value int64 `aper:"valueLB:0,valueUB:255"`
}

// ProcedureCode Value
const (
	ProcedureCodeidAMFConfigurationUpdate                int64 = 0
	ProcedureCodeidAMFStatusIndication                   int64 = 1
	ProcedureCodeidCellTrafficTrace                      int64 = 2
	ProcedureCodeidDeactivateTrace                       int64 = 3
	ProcedureCodeidDownlinkNASTransport                  int64 = 4
	ProcedureCodeidDownlinkNonUEAssociatedNRPPaTransport int64 = 5
	ProcedureCodeidDownlinkRANConfigurationTransfer      int64 = 6
	ProcedureCodeidDownlinkRANStatusTransfer             int64 = 7
	ProcedureCodeidDownlinkUEAssociatedNRPPaTransport    int64 = 8
	ProcedureCodeidErrorIndication                       int64 = 9
	ProcedureCodeidHandoverCancel                        int64 = 10
	ProcedureCodeidHandoverNotification                  int64 = 11
	ProcedureCodeidHandoverPreparation                   int64 = 12
	ProcedureCodeidHandoverResourceAllocation            int64 = 13
	ProcedureCodeidInitialContextSetup                   int64 = 14
	ProcedureCodeidInitialUEMessage                      int64 = 15
	ProcedureCodeidLocationReportingControl              int64 = 16
	ProcedureCodeidLocationReportingFailureIndication    int64 = 17
	ProcedureCodeidLocationReport                        int64 = 18
	ProcedureCodeidNASNonDeliveryIndication              int64 = 19
	ProcedureCodeidNGReset                               int64 = 20
	ProcedureCodeidNGSetup                               int64 = 21
	ProcedureCodeidOverloadStart                         int64 = 22
	ProcedureCodeidOverloadStop                          int64 = 23
	ProcedureCodeidPaging                                int64 = 24
	ProcedureCodeidPathSwitchRequest                     int64 = 25
	ProcedureCodeidPDUSessionResourceModify              int64 = 26
	ProcedureCodeidPDUSessionResourceModifyIndication    int64 = 27
	ProcedureCodeidPDUSessionResourceRelease             int64 = 28
	ProcedureCodeidPDUSessionResourceSetup               int64 = 29
	ProcedureCodeidPDUSessionResourceNotify              int64 = 30
	ProcedureCodeidPrivateMessage                        int64 = 31
	ProcedureCodeidPWSCancel                             int64 = 32
	ProcedureCodeidPWSFailureIndication                  int64 = 33
	ProcedureCodeidPWSRestartIndication                  int64 = 34
	ProcedureCodeidRANConfigurationUpdate                int64 = 35
	ProcedureCodeidRerouteNASRequest                     int64 = 36
	ProcedureCodeidRRCInactiveTransitionReport           int64 = 37
	ProcedureCodeidTraceFailureIndication                int64 = 38
	ProcedureCodeidTraceStart                            int64 = 39
	ProcedureCodeidUECapabilityInfoIndication            int64 = 40
	ProcedureCodeidUEContextModification                 int64 = 41
	ProcedureCodeidUEContextRelease                      int64 = 42
	ProcedureCodeidUEContextReleaseRequest               int64 = 43
	ProcedureCodeidUERadioCapabilityCheck                int64 = 44
	ProcedureCodeidUETNLABindingRelease                  int64 = 45
	ProcedureCodeidUplinkNASTransport                    int64 = 46
	ProcedureCodeidUplinkNonUEAssociatedNRPPaTransport   int64 = 47
	ProcedureCodeidUplinkRANConfigurationTransfer        int64 = 48
	ProcedureCodeidUplinkRANStatusTransfer               int64 = 49
	ProcedureCodeidUplinkUEAssociatedNRPPaTransport      int64 = 50
	ProcedureCodeidWriteReplaceWarning                   int64 = 51
)
