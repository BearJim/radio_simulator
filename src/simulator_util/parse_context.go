package simulator_util

import (
	"fmt"
	"radio_simulator/lib/ngap/ngapConvert"
	"radio_simulator/lib/openapi/models"
	"radio_simulator/src/factory"
	"radio_simulator/src/logger"
	"radio_simulator/src/simulator_context"
	"radio_simulator/src/simulator_nas/nas_security"
	"radio_simulator/src/ue_factory"
	"strconv"
)

func ParseRanContext() {
	config := factory.SimConfig
	self := simulator_context.Simulator_Self()
	self.DefaultRanUri = config.RanInfo[0].RanSctpUri
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
			ran.DefaultTAC = tac
			ran.SupportTAList[tac] = plmnList
		}
	}

}
func ParseUeData(configDirPath string, fileList []string) {
	self := simulator_context.Simulator_Self()
	for _, ueInfoFile := range fileList {
		fileName := configDirPath + ueInfoFile
		ue := ue_factory.InitUeContextFactory(fileName)
		ue.IntAlg = nas_security.AlogMaps[ue.IntegrityAlgOrig]
		ue.EncAlg = nas_security.AlogMaps[ue.CipheringAlgOrig]

		self.UeContextPool[ue.Supi] = ue
	}
}
func InitUeToDB() {
	self := simulator_context.Simulator_Self()
	for supi, ue := range self.UeContextPool {
		amDate := models.AccessAndMobilitySubscriptionData{
			Gpsis: ue.Gpsis,
			Nssai: &ue.Nssai,
		}
		amPolicy := models.AmPolicyData{
			SubscCats: ue.SubscCats,
		}
		auths := ue.AuthData
		authsSubs := models.AuthenticationSubscription{
			AuthenticationMethod:          models.AuthMethod(auths.AuthMethod),
			AuthenticationManagementField: auths.AMF,
			PermanentKey: &models.PermanentKey{
				PermanentKeyValue: auths.K,
			},
			SequenceNumber: auths.SQN,
		}
		if auths.Opc != "" {
			authsSubs.Opc = &models.Opc{
				OpcValue: auths.Opc,
			}
		} else if auths.Op != "" {
			authsSubs.Milenage = &models.Milenage{
				Op: &models.Op{
					OpValue: auths.Op,
				},
			}
		} else {
			logger.UtilLog.Errorf("Ue[%s] need Op or OpCode", ue.Supi)
		}
		InsertAuthSubscriptionToMongoDB(supi, authsSubs)
		InsertAccessAndMobilitySubscriptionDataToMongoDB(supi, amDate, ue.ServingPlmnId)
		InsertSmfSelectionSubscriptionDataToMongoDB(supi, ue.SmfSelData, ue.ServingPlmnId)
		InsertAmPolicyDataToMongoDB(supi, amPolicy)
	}
}

func ClearDB() {
	self := simulator_context.Simulator_Self()
	for supi, ue := range self.UeContextPool {
		logger.UtilLog.Infof("Del UE[%s] Info in DB", supi)
		DelAccessAndMobilitySubscriptionDataFromMongoDB(supi, ue.ServingPlmnId)
		DelAmPolicyDataFromMongoDB(supi)
		DelAuthSubscriptionToMongoDB(supi)
		DelSmfSelectionSubscriptionDataFromMongoDB(supi, ue.ServingPlmnId)
	}
}

func TACConfigToHexString(intString string) (hexString string) {
	tmp, _ := strconv.ParseUint(intString, 10, 32)
	hexString = fmt.Sprintf("%06x", tmp)
	return
}
