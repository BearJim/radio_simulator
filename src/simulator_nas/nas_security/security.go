package nas_security

import (
	"fmt"
	"reflect"

	"github.com/jay16213/radio_simulator/src/simulator_context"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/security"
)

func NASEncode(ue *simulator_context.UeContext, msg *nas.Message, securityContextAvailable bool, newSecurityContext bool) (payload []byte, err error) {
	if ue == nil {
		err = fmt.Errorf("ue is nil")
		return
	}
	if msg == nil {
		err = fmt.Errorf("Nas Message is empty")
		return
	}
	if !securityContextAvailable {
		return msg.PlainNasEncode()
	} else {
		if newSecurityContext {
			ue.ULCount.Set(0, 0)
			ue.DLCount.Set(0, 0)
		}

		sequenceNumber := ue.ULCount.SQN()
		payload, err = msg.PlainNasEncode()
		if err != nil {
			return
		}

		// TODO: Support for ue has nas connection in both accessType
		if err = security.NASEncrypt(ue.EncAlg, ue.KnasEnc, ue.ULCount.Get(), security.Bearer3GPP,
			security.DirectionUplink, payload); err != nil {
			return
		}
		// add sequece number
		payload = append([]byte{sequenceNumber}, payload[:]...)
		mac32 := make([]byte, 4)

		mac32, err = security.NASMacCalculate(ue.IntAlg, ue.KnasInt, ue.ULCount.Get(), security.Bearer3GPP, security.DirectionUplink, payload)
		if err != nil {
			return
		}

		// Add mac value
		payload = append(mac32, payload[:]...)
		// Add EPD and Security Type
		msgSecurityHeader := []byte{msg.SecurityHeader.ProtocolDiscriminator, msg.SecurityHeader.SecurityHeaderType}
		payload = append(msgSecurityHeader, payload[:]...)

		// Increase UL Count
		ue.ULCount.AddOne()
	}
	return
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
	} else if ue.IntAlg == security.AlgIntegrity128NIA0 {
		fmt.Println("decode payload is ", payload)
		// remove header
		payload = payload[3:]

		if err = security.NASEncrypt(ue.EncAlg, ue.KnasEnc, ue.DLCount.Get(), security.Bearer3GPP,
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
				ue.EncAlg = command.SelectedNASSecurityAlgorithms.GetTypeOfCipheringAlgorithm()
				ue.IntAlg = command.SelectedNASSecurityAlgorithms.GetTypeOfIntegrityProtectionAlgorithm()
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

		mac32, err := security.NASMacCalculate(ue.IntAlg, ue.KnasInt, ue.DLCount.Get(), security.Bearer3GPP,
			security.DirectionDownlink, payload)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(mac32, receivedMac32) {
			fmt.Printf("NAS MAC verification failed(0x%x != 0x%x)", mac32, receivedMac32)
		} else {
			fmt.Printf("cmac value: 0x%x\n", mac32)
		}

		// remove sequece Number
		payload = payload[1:]

		// TODO: Support for ue has nas connection in both accessType
		if securityHeaderType != nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext &&
			securityHeaderType != nas.SecurityHeaderTypeIntegrityProtected {
			if err = security.NASEncrypt(ue.EncAlg, ue.KnasEnc, ue.DLCount.Get(), security.Bearer3GPP,
				security.DirectionDownlink, payload); err != nil {
				return nil, err
			}
		}
	}
	err = msg.PlainNasDecode(&payload)
	return
}
