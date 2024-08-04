package logger

import (
	"Ultra-learn/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init() {
	var err error
	zapConfig := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.StacktraceKey = ""
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig = encoderConfig
	zapConfig.Development = config.IS_DEVELOP_MODE
	zapConfig.Encoding = "json"
	zapConfig.InitialFields = map[string]interface{}{"idtx": "999"}
	//zapConfig.OutputPaths = []string{"stdout", config.APP_LOG_FOLDER + "app_log.log"}
	//zapConfig.ErrorOutputPaths = []string{"stderr"}
	log, err = zapConfig.Build(zap.AddCallerSkip(1))
	defer log.Sync()
	//log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

}

func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}
func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}
func Error(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}
