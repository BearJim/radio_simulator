package path_util

import (
	"radio_simulator/lib/path_util/logger"
	"testing"
)

func TestFree5gcPath(t *testing.T) {
	logger.PathLog.Infoln(ModulePath("gofree5gc/abcdef/abcdef.pem"))
}
