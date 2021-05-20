package nas_security

import (
	"fmt"
	"reflect"

	"github.com/jay16213/radio_simulator/pkg/logger"
	"github.com/jay16213/radio_simulator/pkg/simulator_context"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/security"
)

func NASEncode(ue *simulator_context.UeContext, msg *nas.Message, securityContextAvailable bool, newSecurityContext bool) ([]byte, error) {
	if ue == nil {
		return nil, fmt.Errorf("ue is nil")
	}
	if msg == nil {
		return nil, fmt.Errorf("Nas Message is empty")
	}
	if !securityContextAvailable {
		return msg.PlainNasEncode()
	} else {
		if newSecurityContext {
			ue.ULCount.Set(0, 0)
			ue.DLCount.Set(0, 0)
		}

		sequenceNumber := ue.ULCount.SQN()
		payload, err := msg.PlainNasEncode()
		if err != nil {
			return nil, err
		}

		// TODO: Support for ue has nas connection in both accessType
		logger.NASLog.Debugf("Encrypt NAS message (algorithm: %+v, ULCount: 0x%0x)", ue.CipheringAlg, ue.ULCount.ToUint32())
		logger.NASLog.Debugf("NAS ciphering key: %0x", ue.KnasEnc)
		if msg.SecurityHeaderType != nas.SecurityHeaderTypeIntegrityProtected {
			if err = security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.ULCount.ToUint32(), security.Bearer3GPP,
				security.DirectionUplink, payload); err != nil {
				return nil, err
			}
		}
		// add sequece number
		payload = append([]byte{sequenceNumber}, payload[:]...)

		mac32, err := security.NASMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.ULCount.ToUint32(), security.Bearer3GPP, security.DirectionUplink, payload)
		if err != nil {
			return nil, err
		}

		// Add mac value
		payload = append(mac32, payload[:]...)
		// Add EPD and Security Type
		msgSecurityHeader := []byte{msg.SecurityHeader.ProtocolDiscriminator, msg.SecurityHeader.SecurityHeaderType}
		payload = append(msgSecurityHeader, payload[:]...)

		// Increase UL Count
		ue.ULCount.AddOne()
		return payload, nil
	}
}

func NASDecode(ue *simulator_context.UeContext, securityHeaderType uint8, payload []byte) (msg *nas.Message, err error) {
	if ue == nil {
		err = fmt.Errorf("ue is nil")
		return
	}
	if payload == nil {
		err = fmt.Errorf("Nas payload is empty")
		return
	}

	msg = new(nas.Message)

	if securityHeaderType == nas.SecurityHeaderTypePlainNas {
		err = msg.PlainNasDecode(&payload)
		return
	} else if ue.IntegrityAlg == security.AlgIntegrity128NIA0 {
		// remove header
		payload = payload[3:]

		if err = security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.DLCount.ToUint32(), security.Bearer3GPP,
			security.DirectionDownlink, payload); err != nil {
			return nil, err
		}

		err = msg.PlainNasDecode(&payload)
		return
	} else {
		// security mode command
		if securityHeaderType == nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext {
			ue.DLCount.Set(0, 0)

			plainNas := payload[7:]
			if err := msg.PlainNasDecode(&plainNas); err != nil {
				return nil, err
			}
			if command := msg.GmmMessage.SecurityModeCommand; command != nil {
				ue.CipheringAlg = command.SelectedNASSecurityAlgorithms.GetTypeOfCipheringAlgorithm()
				ue.IntegrityAlg = command.SelectedNASSecurityAlgorithms.GetTypeOfIntegrityProtectionAlgorithm()
				ue.DerivateAlgKey()
			} else {
				return nil, fmt.Errorf("Integrity Protected With New 5G Nas Security is not Security command")
			}
		}

		securityHeader := payload[0:6]
		sequenceNumber := payload[6]
		receivedMac32 := securityHeader[2:]
		// remove security Header except for sequece Number
		payload = payload[6:]

		// Caculate dl count
		if ue.DLCount.SQN() > sequenceNumber {
			ue.DLCount.SetOverflow(ue.DLCount.Overflow() + 1)
		}
		ue.DLCount.SetSQN(sequenceNumber)

		logger.NASLog.Debugf("Calculate NAS MAC (algorithm: %+v, DLCount: 0x%0x)", ue.IntegrityAlg, ue.DLCount.ToUint32())
		logger.NASLog.Debugf("NAS integrity key: %0x", ue.KnasInt)
		mac32, err := security.NASMacCalculate(ue.IntegrityAlg, ue.KnasInt, ue.DLCount.ToUint32(), security.Bearer3GPP,
			security.DirectionDownlink, payload)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(mac32, receivedMac32) {
			logger.NASLog.Warnf("NAS MAC verification failed(0x%x != 0x%x)", mac32, receivedMac32)
		} else {
			logger.NASLog.Debugf("cmac value: 0x%x\n", mac32)
		}

		// remove sequece Number
		payload = payload[1:]

		// TODO: Support for ue has nas connection in both accessType
		if securityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext &&
			securityHeaderType != nas.SecurityHeaderTypeIntegrityProtected {
			logger.NASLog.Debugf("Decrypt NAS message (algorithm: %+v, DLCount: 0x%0x)", ue.CipheringAlg, ue.DLCount.ToUint32())
			logger.NASLog.Debugf("NAS ciphering key: %0x", ue.KnasEnc)
			if err = security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.DLCount.ToUint32(), security.Bearer3GPP,
				security.DirectionDownlink, payload); err != nil {
				return nil, err
			}
		}
	}
	err = msg.PlainNasDecode(&payload)
	return
}
