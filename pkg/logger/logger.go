package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

var AppLog *logrus.Entry
var UtilLog *logrus.Entry
var NgapLog *logrus.Entry
var GtpLog *logrus.Entry
var NASLog *logrus.Entry
var ApiLog *logrus.Entry
var ContextLog *logrus.Entry

// var zapCfg zap.Config
// var AppLog *zap.Logger

func init() {
	// zapCfg = zap.NewProductionConfig()
	// zapCfg.DisableCaller = true
	// zapCfg.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
	// 	pae.AppendString(t.Format(time.RFC3339))
	// }
	// zapCfg.OutputPaths = append(zapCfg.OutputPaths, "./log/ran.log")
	// zapLog, err := zapCfg.Build()
	// defer zapLog.Sync()
	// if err != nil {
	// 	panic(err)
	// }
	// AppLog = zapLog.With(zap.String("app", "ran"))

	log = logrus.New()
	log.SetReportCaller(true)

	log.Formatter = &logrus.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			orgFilename, _ := os.Getwd()
			repopath := orgFilename
			repopath = strings.Replace(repopath, "/bin", "", 1)
			filename := strings.Replace(f.File, repopath, "", -1)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	AppLog = log.WithFields(logrus.Fields{"RAN": "App"})
	NgapLog = log.WithFields(logrus.Fields{"RAN": "NGAP"})
	GtpLog = log.WithFields(logrus.Fields{"RAN": "GTP"})
	UtilLog = log.WithFields(logrus.Fields{"RAN": "Util"})
	NASLog = log.WithFields(logrus.Fields{"RAN": "NAS"})
	ApiLog = log.WithFields(logrus.Fields{"RAN": "API"})
	ContextLog = log.WithFields(logrus.Fields{"RAN": "Context"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
