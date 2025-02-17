package logging

import (
	"gango/utils"
	"path/filepath"
)

type Logging struct {
}

func (l Logging) WriteFolder(dir string) error {
	return utils.WriteFile(dir, filepath.Join(l.FilePath(), l.FileName()), loggingFile)
}

func (l Logging) FilePath() string {
	return "/src/lib/logging"
}

func (l Logging) FileName() string {
	return "logging.go"
}

var loggingFile = `
package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerFactory interface {
	New() *zap.SugaredLogger
}

var (
	loggerFactory LoggerFactory
	loggers       = map[string]*zap.SugaredLogger{}
)

func SetFactory(factory LoggerFactory) {
	loggerFactory = factory
}

func SetDefaultFactory() {
	loggerFactory = LoggerFactory(ZapLoggerFactory{
		Config: LogConfiguration{
			Path:        "./service.log",
			Level:       -1,
			Development: false,
		},
	})
}

func GetLogger(name string) *zap.SugaredLogger {
	if loggerFactory == nil {
		SetDefaultFactory()
	}
	if _, ok := loggers[name]; !ok {
		loggers[name] = loggerFactory.New()
	}
	return loggers[name]
}

type ZapLoggerFactory struct {
	Config LogConfiguration
}

func (d ZapLoggerFactory) New() *zap.SugaredLogger {
	logger := NewZapSugarLogger(d.Config)
	defer logger.Sync()
	return logger
}

type LogConfiguration struct {
	Path        string
	Level       int
	Development bool
}

func NewZapSugarLogger(config LogConfiguration) *zap.SugaredLogger {
	path := config.Path
	zapBuilder := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(config.Level)),
		Development: config.Development,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			TimeKey:     "timestamp",
			EncodeLevel: zapcore.LowercaseColorLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
		OutputPaths:      []string{path, "stdout"},
		ErrorOutputPaths: []string{path, "stderr"},
	}
	logger, _ := zapBuilder.Build()
	return logger.Sugar()
}
`
