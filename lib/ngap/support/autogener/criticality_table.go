package ngap_autogener

import (
	"radio_simulator/lib/ngap/ngapType"
)

var presentMap = map[int]string{
	ngapType.NGAPPDUPresentInitiatingMessage:   "InitiatingMessage",
	ngapType.NGAPPDUPresentSuccessfulOutcome:   "SuccessfulOutcome",
	ngapType.NGAPPDUPresentUnsuccessfulOutcome: "UnsuccessfulOutcome",
}

var initiatingMap = map[int][]int{
	ngapType.InitiatingMessagePresentNothing:                               {},
	ngapType.InitiatingMessagePresentAMFConfigurationUpdate:                {1, 2, 2, 0, 2, 0, 0, 0}, // 1 2 4 reject but optional
	ngapType.InitiatingMessagePresentHandoverCancel:                        {1, 1, 1, 0},
	ngapType.InitiatingMessagePresentHandoverRequired:                      {1, 1, 1, 1, 0, 1, 0, 1, 1},
	ngapType.InitiatingMessagePresentHandoverRequest:                       {1, 1, 1, 0, 1, 0, 1, 1, 2, 2, 1, 1, 0, 0, 1, 0, 0, 0, 1},    // 8 9 reject but optional
	ngapType.InitiatingMessagePresentInitialContextSetupRequest:            {1, 1, 1, 2, 1, 0, 1, 2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 2, 0, 0}, // 3 7 17 reject but optional
	ngapType.InitiatingMessagePresentNGReset:                               {1, 0, 1},
	ngapType.InitiatingMessagePresentNGSetupRequest:                        {1, 1, 0, 1, 0},
	ngapType.InitiatingMessagePresentPathSwitchRequest:                     {1, 1, 1, 0, 0, 1, 0},
	ngapType.InitiatingMessagePresentPDUSessionResourceModifyRequest:       {1, 1, 1, 0, 1},
	ngapType.InitiatingMessagePresentPDUSessionResourceModifyIndication:    {1, 1, 1, 1},
	ngapType.InitiatingMessagePresentPDUSessionResourceReleaseCommand:      {1, 1, 1, 0, 0, 1},
	ngapType.InitiatingMessagePresentPDUSessionResourceSetupRequest:        {1, 1, 1, 0, 2, 1},                   // 4 reject but optional
	ngapType.InitiatingMessagePresentPWSCancelRequest:                      {1, 1, 1, 0, 2},                      // 4 reject but optional
	ngapType.InitiatingMessagePresentRANConfigurationUpdate:                {1, 0, 2, 0},                         // 2 reject but optional
	ngapType.InitiatingMessagePresentUEContextModificationRequest:          {1, 1, 1, 0, 2, 0, 0, 2, 0, 2, 0, 0}, // 4 7 9 reject but optional
	ngapType.InitiatingMessagePresentUEContextReleaseCommand:               {1, 1, 0},
	ngapType.InitiatingMessagePresentUERadioCapabilityCheckRequest:         {1, 1, 1, 0},
	ngapType.InitiatingMessagePresentWriteReplaceWarningRequest:            {1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 2, 0}, // 10 reject but optional
	ngapType.InitiatingMessagePresentAMFStatusIndication:                   {0, 1},
	ngapType.InitiatingMessagePresentCellTrafficTrace:                      {0, 1, 1, 0, 0, 0},
	ngapType.InitiatingMessagePresentDeactivateTrace:                       {0, 1, 1, 0},
	ngapType.InitiatingMessagePresentDownlinkNASTransport:                  {0, 1, 1, 2, 0, 1, 0, 0, 0, 2}, // 3 9 reject but optional
	ngapType.InitiatingMessagePresentDownlinkNonUEAssociatedNRPPaTransport: {0, 1, 1},
	ngapType.InitiatingMessagePresentDownlinkRANConfigurationTransfer:      {0, 0},
	ngapType.InitiatingMessagePresentDownlinkRANStatusTransfer:             {0, 1, 1, 1},
	ngapType.InitiatingMessagePresentDownlinkUEAssociatedNRPPaTransport:    {0, 1, 1, 1, 1},
	ngapType.InitiatingMessagePresentErrorIndication:                       {0, 0, 0, 0, 0},
	ngapType.InitiatingMessagePresentHandoverNotify:                        {0, 1, 1, 0},
	ngapType.InitiatingMessagePresentInitialUEMessage:                      {0, 1, 1, 1, 0, 2, 0, 0, 2}, // 5 8 reject but optional
	ngapType.InitiatingMessagePresentLocationReport:                        {0, 1, 1, 0, 0, 0},
	ngapType.InitiatingMessagePresentLocationReportingControl:              {0, 1, 1, 0},
	ngapType.InitiatingMessagePresentLocationReportingFailureIndication:    {0, 1, 1, 0},
	ngapType.InitiatingMessagePresentNASNonDeliveryIndication:              {0, 1, 1, 0, 0},
	ngapType.InitiatingMessagePresentOverloadStart:                         {0, 2, 0, 0}, // 1 reject but optional
	ngapType.InitiatingMessagePresentOverloadStop:                          {1},
	ngapType.InitiatingMessagePresentPaging:                                {0, 0, 0, 0, 0, 0, 0, 0},
	ngapType.InitiatingMessagePresentPDUSessionResourceNotify:              {0, 1, 1, 2, 0, 0}, // 3 reject but optional
	ngapType.InitiatingMessagePresentPrivateMessage:                        {0},
	ngapType.InitiatingMessagePresentPWSFailureIndication:                  {0, 1, 1},
	ngapType.InitiatingMessagePresentPWSRestartIndication:                  {0, 1, 1, 1, 2},    // 4 reject but optional
	ngapType.InitiatingMessagePresentRerouteNASRequest:                     {1, 1, 0, 1, 1, 2}, // 5 reject but optional
	ngapType.InitiatingMessagePresentRRCInactiveTransitionReport:           {0, 1, 1, 0, 0},
	ngapType.InitiatingMessagePresentTraceFailureIndication:                {0, 1, 1, 0, 0},
	ngapType.InitiatingMessagePresentTraceStart:                            {0, 1, 1, 0},
	ngapType.InitiatingMessagePresentUEContextReleaseRequest:               {0, 1, 1, 2, 0}, // 3 reject but optional
	ngapType.InitiatingMessagePresentUERadioCapabilityInfoIndication:       {0, 1, 1, 0, 0},
	ngapType.InitiatingMessagePresentUETNLABindingReleaseRequest:           {0, 1, 1},
	ngapType.InitiatingMessagePresentUplinkNASTransport:                    {0, 1, 1, 1, 0},
	ngapType.InitiatingMessagePresentUplinkNonUEAssociatedNRPPaTransport:   {0, 1, 1},
	ngapType.InitiatingMessagePresentUplinkRANConfigurationTransfer:        {0, 0},
	ngapType.InitiatingMessagePresentUplinkRANStatusTransfer:               {0, 1, 1, 1},
	ngapType.InitiatingMessagePresentUplinkUEAssociatedNRPPaTransport:      {0, 1, 1, 1, 1},
}

var successfulMap = map[int][]int{
	ngapType.SuccessfulOutcomePresentNothing:                           {},
	ngapType.SuccessfulOutcomePresentAMFConfigurationUpdateAcknowledge: {1, 0, 0, 0},
	ngapType.SuccessfulOutcomePresentHandoverCancelAcknowledge:         {1, 0, 0, 0},
	ngapType.SuccessfulOutcomePresentHandoverCommand:                   {1, 1, 1, 1, 1, 0, 0, 1, 0},
	ngapType.SuccessfulOutcomePresentHandoverRequestAcknowledge:        {1, 0, 0, 0, 0, 1, 0},
	ngapType.SuccessfulOutcomePresentInitialContextSetupResponse:       {1, 0, 0, 0, 0, 0},
	ngapType.SuccessfulOutcomePresentNGResetAcknowledge:                {1, 0, 0},
	ngapType.SuccessfulOutcomePresentNGSetupResponse:                   {1, 1, 1, 0, 1, 0},
	ngapType.SuccessfulOutcomePresentPathSwitchRequestAcknowledge:      {1, 0, 0, 2, 1, 0, 0, 0, 1, 0, 0, 0}, // 3 reject but optional
	ngapType.SuccessfulOutcomePresentPDUSessionResourceModifyResponse:  {1, 0, 0, 0, 0, 0, 0},
	ngapType.SuccessfulOutcomePresentPDUSessionResourceModifyConfirm:   {1, 0, 0, 0, 0, 0},
	ngapType.SuccessfulOutcomePresentPDUSessionResourceReleaseResponse: {1, 0, 0, 0, 0, 0},
	ngapType.SuccessfulOutcomePresentPDUSessionResourceSetupResponse:   {1, 0, 0, 0, 0, 0},
	ngapType.SuccessfulOutcomePresentPWSCancelResponse:                 {1, 1, 1, 0, 0},
	ngapType.SuccessfulOutcomePresentRANConfigurationUpdateAcknowledge: {1, 0},
	ngapType.SuccessfulOutcomePresentUEContextModificationResponse:     {1, 0, 0, 0, 0, 0},
	ngapType.SuccessfulOutcomePresentUEContextReleaseComplete:          {1, 0, 0, 0, 0, 2, 0}, // 5 reject but optional
	ngapType.SuccessfulOutcomePresentUERadioCapabilityCheckResponse:    {1, 0, 0, 1, 0},
	ngapType.SuccessfulOutcomePresentWriteReplaceWarningResponse:       {1, 1, 1, 0, 0},
}

var unsuccessfulMap = map[int][]int{
	ngapType.UnsuccessfulOutcomePresentNothing:                       {},
	ngapType.UnsuccessfulOutcomePresentAMFConfigurationUpdateFailure: {1, 0, 0, 0},
	ngapType.UnsuccessfulOutcomePresentHandoverPreparationFailure:    {1, 0, 0, 0, 0},
	ngapType.UnsuccessfulOutcomePresentHandoverFailure:               {1, 0, 0, 0},
	ngapType.UnsuccessfulOutcomePresentInitialContextSetupFailure:    {1, 0, 0, 0, 0, 0},
	ngapType.UnsuccessfulOutcomePresentNGSetupFailure:                {1, 0, 0, 0},
	ngapType.UnsuccessfulOutcomePresentPathSwitchRequestFailure:      {1, 0, 0, 0, 0},
	ngapType.UnsuccessfulOutcomePresentRANConfigurationUpdateFailure: {1, 0, 0, 0},
	ngapType.UnsuccessfulOutcomePresentUEContextModificationFailure:  {1, 0, 0, 0, 0},
}

var criticalityMap = map[int]map[int][]int{
	ngapType.NGAPPDUPresentInitiatingMessage:   initiatingMap,
	ngapType.NGAPPDUPresentSuccessfulOutcome:   successfulMap,
	ngapType.NGAPPDUPresentUnsuccessfulOutcome: unsuccessfulMap,
}

var criticalityToString = map[int]string{
	0: "Ignore",
	1: "Reject",
	2: "Reject",
}
