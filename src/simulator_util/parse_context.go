package simulator_util

import (
	"fmt"
	"radio_simulator/lib/ngap/ngapConvert"
	"radio_simulator/src/factory"
	"radio_simulator/src/simulator_context"
	"strconv"
)

func ParseRanContext() {
	config := factory.SimConfig
	self := simulator_context.Simulator_Self()
	for _, ranInfo := range config.RanInfo {
		plmnId := ngapConvert.PlmnIdToNgap(ranInfo.GnbId.PlmnId)
		ran := self.AddRanContext(ranInfo.AmfUri, ranInfo.RanSctpUri, ranInfo.RanName, plmnId, ranInfo.GnbId.Value, ranInfo.GnbId.BitLength)
		for _, supportItem := range ranInfo.SupportTAList {
			plmnList := []simulator_context.PlmnSupportItem{}
			for _, item := range supportItem.Plmnlist {
				plmnItem := simulator_context.PlmnSupportItem{}
				plmnItem.PlmnId = ngapConvert.PlmnIdToNgap(item.PlmnId)
				for _, snssai := range item.SNssaiList {
					sNssaiNgap := ngapConvert.SNssaiToNgap(snssai)
					plmnItem.SNssaiList = append(plmnItem.SNssaiList, sNssaiNgap)
				}
				plmnList = append(plmnList, plmnItem)
			}
			tac := TACConfigToHexString(supportItem.Tac)
			ran.SupportTAList[tac] = plmnList
		}
	}
}
func TACConfigToHexString(intString string) (hexString string) {
	tmp, _ := strconv.ParseUint(intString, 10, 32)
	hexString = fmt.Sprintf("%06x", tmp)
	return
}
