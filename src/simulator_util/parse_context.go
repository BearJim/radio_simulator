package simulator_util

import (
	"fmt"
	"radio_simulator/lib/ngap/ngapConvert"
	"radio_simulator/lib/openapi/models"
	"radio_simulator/src/factory"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/ue_factory"
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
func ParseUeData(configDirPath string, fileList []string) {
	self := simulator_context.Simulator_Self()
	for _, ueInfoFile := range fileList {
		fileName := configDirPath + ueInfoFile
		config := ue_factory.InitUeConfigFactory(fileName)
		ueInfo := simulator_context.UeDBInfo{}
		ueInfo.AmDate = models.AccessAndMobilitySubscriptionData{
			Gpsis: config.Gpsis,
			Nssai: &config.Nssai,
		}
		ueInfo.SmfSelData = config.SmfSelData
		ueInfo.PlmnId = config.ServingPlmnId
		ueInfo.AmPolicy.SubscCats = config.SubscCats
		auths := config.AuthData
		ueInfo.AuthsSubs = models.AuthenticationSubscription{
			AuthenticationMethod:          models.AuthMethod(auths.AuthMethod),
			AuthenticationManagementField: auths.AMF,
			PermanentKey: &models.PermanentKey{
				PermanentKeyValue: auths.K,
			},
			SequenceNumber: auths.SQN,
		}
		if auths.Opc != "" {
			ueInfo.AuthsSubs.Opc = &models.Opc{
				OpcValue: auths.Opc,
			}
		} else if auths.Op != "" {
			ueInfo.AuthsSubs.Milenage = &models.Milenage{
				Op: &models.Op{
					OpValue: auths.Op,
				},
			}
		} else {
			logger.UtilLog.Errorf("Ue[%s] need Op or OpCode", config.Supi)
		}
		self.UeContextPool[config.Supi] = ueInfo
	}
}
func InitUeToDB() {
	self := simulator_context.Simulator_Self()
	for supi, info := range self.UeContextPool {
		InsertAuthSubscriptionToMongoDB(supi, info.AuthsSubs)
		InsertAccessAndMobilitySubscriptionDataToMongoDB(supi, info.AmDate, info.PlmnId)
		InsertSmfSelectionSubscriptionDataToMongoDB(supi, info.SmfSelData, info.PlmnId)
		InsertAmPolicyDataToMongoDB(supi, info.AmPolicy)
	}
}

func ClearDB() {
	self := simulator_context.Simulator_Self()
	for supi, info := range self.UeContextPool {
		logger.UtilLog.Infof("Del UE[%s] Info in DB", supi)
		DelAccessAndMobilitySubscriptionDataFromMongoDB(supi, info.PlmnId)
		DelAmPolicyDataFromMongoDB(supi)
		DelAuthSubscriptionToMongoDB(supi)
		DelSmfSelectionSubscriptionDataFromMongoDB(supi, info.PlmnId)
	}
}

func TACConfigToHexString(intString string) (hexString string) {
	tmp, _ := strconv.ParseUint(intString, 10, 32)
	hexString = fmt.Sprintf("%06x", tmp)
	return
}
