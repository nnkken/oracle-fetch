package utils

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/spf13/cobra"
)

const (
	CmdLogLevel   = "log-level"
	CmdLogFormat  = "log-format"
	CmdLogOutputs = "log-outputs"

	DefaultLogLevel  = "info"
	DefaultLogFormat = "console"
)

var (
	DefaultLogOutputs = []string{"stderr"}

	cfg = zap.NewProductionConfig()

	loggersLock = &sync.Mutex{}
	loggers     = make(map[string]*zap.SugaredLogger)
)

func ConfigCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String(CmdLogLevel, DefaultLogLevel, "logging level (debug | info | warn | error | dpanic | panic | fatal)")
	cmd.PersistentFlags().String(CmdLogFormat, DefaultLogFormat, "logging format (json | console)")
	cmd.PersistentFlags().StringArray(CmdLogOutputs, DefaultLogOutputs, "logging outputs (stdout | stderr | /somewhere/to/some/file)")
}

func SetupLoggerFromCmdArgs(cmd *cobra.Command) error {
	cmdLevel, err := cmd.Flags().GetString(CmdLogLevel)
	if err != nil {
		return err
	}
	cmdFormat, err := cmd.Flags().GetString(CmdLogFormat)
	if err != nil {
		return err
	}
	cmdOutputs, err := cmd.Flags().GetStringArray(CmdLogOutputs)
	if err != nil {
		return err
	}
	var level zapcore.Level
	err = level.UnmarshalText([]byte(cmdLevel))
	if err != nil {
		return err
	}
	SetupLoggerConfig(level, cmdOutputs, cmdFormat)
	return nil
}

func SetupLoggerConfig(level zapcore.Level, cmdOutputs []string, cmdFormat string) {
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.OutputPaths = cmdOutputs
	cfg.ErrorOutputPaths = cmdOutputs
	cfg.Encoding = cmdFormat
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
}

func GetLogger(module string) *zap.SugaredLogger {
	loggersLock.Lock()
	defer loggersLock.Unlock()
	subLogger, ok := loggers[module]
	if !ok {
		rawLogger, err := cfg.Build()
		if err != nil {
			// seems we don't have a proper logger config, so can only use the default one and panic
			zap.S().Panicw("failed to build logger", "error", err)
		}
		subLogger = rawLogger.Sugar().Named(module)
		loggers[module] = subLogger
	}
	return subLogger
}
