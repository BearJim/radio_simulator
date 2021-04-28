package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapCfg zap.Config
	zapLog *zap.Logger
)

var (
	AppLog     *zap.SugaredLogger
	NgapLog    *zap.SugaredLogger
	ApiLog     *zap.SugaredLogger
	ContextLog *zap.SugaredLogger
	NASLog     *zap.SugaredLogger
)

func init() {
	zapCfg = zap.NewProductionConfig()
	zapCfg.DisableCaller = true
	zapCfg.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.RFC3339))
	}

	zapCfg.OutputPaths = append(zapCfg.OutputPaths, "./ran.log")
	zapCfg.ErrorOutputPaths = append(zapCfg.ErrorOutputPaths, "./ran.log")
	log, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := log.Sync(); err != nil && !strings.Contains(err.Error(), os.ErrInvalid.Error()) {
			fmt.Println(err)
		}
	}()

	zapLog = log

	AppLog = zapLog.With(zap.String("RAN", "APP")).Sugar()
	NgapLog = zapLog.With(zap.String("RAN", "NGAP")).Sugar()
	ApiLog = zapLog.With(zap.String("RAN", "API")).Sugar()
	ContextLog = zapLog.With(zap.String("RAN", "CTX")).Sugar()
	NASLog = zapLog.With(zap.String("RAN", "NAS")).Sugar()
}

func SetLogLevel(level logrus.Level) {
	zapCfg.Level.SetLevel(zap.DebugLevel)
}

func SetReportCaller(bool bool) {
	// log.SetReportCaller(bool)
}
