package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

type RecommendedRANNodesForPaging struct {
	RecommendedRANNodeList RecommendedRANNodeList
	IEExtensions           *ProtocolExtensionContainerRecommendedRANNodesForPagingExtIEs `aper:"optional"`
}
