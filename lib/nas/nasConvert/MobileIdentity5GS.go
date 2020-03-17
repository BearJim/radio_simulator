package nasConvert

import (
	"encoding/hex"
	"radio_simulator/lib/nas/nasMessage"
	"radio_simulator/lib/nas/nasType"
	"radio_simulator/lib/openapi/models"
	"strconv"
)

func GetTypeOfIdentity(buf byte) uint8 {
	return buf & 0x07
}

// TS 24.501 9.11 .3.4
func GetPlmnIdAndSupiFromSuci(buf []byte) (plmnId, supi string) {

	supiFormat := (buf[0] & 0xf0) >> 4

	var prefix string

	switch supiFormat {
	case 0x00: // imsi
		prefix = "imsi-"
	case 0x01: // Network Specific Identifier
		prefix = "nai-"
	default: // All other values are interpreted as IMSI by this version of the protocol
		prefix = "imsi-"
	}

	// TODO: process Network Specific Identifier

	plmnId = PlmnIDToString(buf[1:4])

	var tmpBytes []byte
	if buf[6] == 0x00 { // protectionSchemeID is "Null scheme"
		for _, rawMsin := range buf[8:] {
			data := ((rawMsin & 0x0f) << 4) | ((rawMsin & 0xf0) >> 4)
			tmpBytes = append(tmpBytes, uint8(data))
		}
	}

	msin := hex.EncodeToString(tmpBytes)
	if msin[len(msin)-1] == 'f' {
		msin = msin[:len(msin)-1] // cut unused digit
	}

	supi = prefix + plmnId + msin
	return
}

// nasType: TS 24.501 9.11.3.4
func GutiToString(buf []byte) (guami models.Guami, guti string) {

	plmnID := PlmnIDToString(buf[1:4])
	amfID := hex.EncodeToString(buf[4:7])
	tmsi5G := hex.EncodeToString(buf[7:])

	guami.PlmnId = new(models.PlmnId)
	guami.PlmnId.Mcc = plmnID[:3]
	guami.PlmnId.Mnc = plmnID[3:]
	guami.AmfId = amfID
	guti = plmnID + amfID + tmsi5G
	return
}

func GutiToNas(guti string) (gutiNas nasType.GUTI5G) {

	gutiNas.SetLen(11)
	gutiNas.SetSpare(0)
	gutiNas.SetTypeOfIdentity(nasMessage.MobileIdentity5GSType5gGuti)

	mcc1, _ := strconv.Atoi(string(guti[0]))
	mcc2, _ := strconv.Atoi(string(guti[1]))
	mcc3, _ := strconv.Atoi(string(guti[2]))
	mnc1, _ := strconv.Atoi(string(guti[3]))
	mnc2, _ := strconv.Atoi(string(guti[4]))
	mnc3 := 0x0f
	amfId := ""
	tmsi := ""
	if len(guti) == 20 {
		mnc3, _ = strconv.Atoi(string(guti[5]))
		amfId = guti[6:12]
		tmsi = guti[12:]
	} else {
		amfId = guti[5:11]
		tmsi = guti[11:]
	}
	gutiNas.SetMCCDigit1(uint8(mcc1))
	gutiNas.SetMCCDigit2(uint8(mcc2))
	gutiNas.SetMCCDigit3(uint8(mcc3))
	gutiNas.SetMNCDigit1(uint8(mnc1))
	gutiNas.SetMNCDigit2(uint8(mnc2))
	gutiNas.SetMNCDigit3(uint8(mnc3))

	amfRegionId, amfSetId, amfPointer := AmfIdToNas(amfId)
	gutiNas.SetAMFRegionID(amfRegionId)
	gutiNas.SetAMFSetID(amfSetId)
	gutiNas.SetAMFPointer(amfPointer)
	tmsiBytes, _ := hex.DecodeString(tmsi)
	copy(gutiNas.Octet[7:11], tmsiBytes[:])
	return
}

// PEI: ^(imei-[0-9]{15}|imeisv-[0-9]{16}|.+)$
func PeiToString(buf []byte) (pei string) {

	var prefix string

	typeOfIdentity := buf[0] & 0x07
	if typeOfIdentity == 0x03 {
		prefix = "imei-"
	} else {
		prefix = "imeisv-"
	}

	oddIndication := (buf[0] & 0x08) >> 3

	digit1 := (buf[0] & 0xf0)

	tmpBytes := []byte{digit1}

	for _, octet := range buf[1:] {
		digitP := octet & 0x0f
		digitP1 := octet & 0xf0

		tmpBytes[len(tmpBytes)-1] += digitP
		tmpBytes = append(tmpBytes, digitP1)
	}

	digitStr := hex.EncodeToString(tmpBytes)
	digitStr = digitStr[:len(digitStr)-1] // remove the last digit

	if oddIndication == 0 { // even digits
		digitStr = digitStr[:len(digitStr)-1] // remove the last digit
	}

	pei = prefix + digitStr
	return
}
