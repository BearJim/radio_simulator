package simulator_util

import (
	"bitbucket.org/free5gc-team/nas"
	"bitbucket.org/free5gc-team/ngap/ngapType"
)

func GetNasPdu(msg *ngapType.DownlinkNASTransport) (m *nas.Message) {
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDNASPDU {
			pkg := []byte(ie.Value.NASPDU.Value)
			m = new(nas.Message)
			err := m.PlainNasDecode(&pkg)
			if err != nil {
				return nil
			}
			return
		}
	}
	return nil
}
