package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type PagingAttemptInformation struct {
	PagingAttemptCount             PagingAttemptCount
	IntendedNumberOfPagingAttempts IntendedNumberOfPagingAttempts
	NextPagingAreaScope            *NextPagingAreaScope                                      `aper:"optional"`
	IEExtensions                   *ProtocolExtensionContainerPagingAttemptInformationExtIEs `aper:"optional"`
}
