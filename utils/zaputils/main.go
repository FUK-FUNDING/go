package zaputils

import (
	"fmt"
	"fuk-funding/go/fp"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

func getColoredLevel(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return color.BlueString("DEBUG")
	case zapcore.InfoLevel:
		return color.GreenString("INFO ")
	case zapcore.WarnLevel:
		return color.YellowString("WARN ")
	case zapcore.ErrorLevel:
		return color.RedString("ERROR")
	default:
		return level.CapitalString()
	}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(color.CyanString(t.Format("15:04:05.000000")))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(getColoredLevel(level))
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	path := caller.TrimmedPath()
	numbers := strings.Split(path, ":")
	if len(numbers) != 2 {
		enc.AppendString(color.MagentaString(path))
		return
	}

	if len(numbers[1]) == 1 {
		path += "  "
	} else if len(numbers[1]) == 2 {
		path += " "
	}

	enc.AppendString(color.MagentaString(path))
}

func customNameEncoder(name string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(color.WhiteString(name))
}

func fileCore(logFile string) zapcore.Core {
	err := fp.StoreFileRecursive(logFile, []byte{})
	if err != nil {
		return nil
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open log file: %v\n", err)
	}

	fileConfig := zap.NewProductionEncoderConfig()
	fileConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileConfig.EncodeCaller = zapcore.FullCallerEncoder
	fileConfig.EncodeName = zapcore.FullNameEncoder

	fileEncoder := zapcore.NewConsoleEncoder(fileConfig)

	return zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zap.DebugLevel)
}

func consoleCore() zapcore.Core {
	// Console encoder config
	consoleConfig := zap.NewProductionEncoderConfig()
	consoleConfig.EncodeTime = customTimeEncoder
	consoleConfig.EncodeLevel = customLevelEncoder
	consoleConfig.EncodeCaller = customCallerEncoder
	consoleConfig.EncodeName = customNameEncoder
	consoleConfig.ConsoleSeparator = " "

	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	return zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
}

func InitLogger() *zap.Logger {
	core := zapcore.NewTee(
		consoleCore(),
		//fileCore(filepath.Join("./logs/general.log")),
	)

	return zap.New(core, zap.AddCaller())
}
