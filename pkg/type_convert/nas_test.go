package type_convert_test

import (
	"testing"

	"github.com/BearJim/radio_simulator/pkg/type_convert"

	"github.com/free5gc/nas/nasType"
	"github.com/stretchr/testify/assert"
)

func TestSupiToMobileId(t *testing.T) {
	supi := "imsi-2089300007487"
	mobileId := type_convert.SupiToMobileId(supi, "20893")
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    12, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
	}
	assert.Equal(t, mobileIdentity5GS, mobileId)
}
