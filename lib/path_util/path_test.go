package path_util

import (
	"testing"

	"github.com/jay16213/radio_simulator/lib/path_util/logger"
)

func TestFree5gcPath(t *testing.T) {
	logger.PathLog.Infoln(ModulePath("gofree5gc/abcdef/abcdef.pem"))
}
