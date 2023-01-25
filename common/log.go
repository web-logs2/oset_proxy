//
// File: log.go
// Created by Dizzrt on 2023/01/21.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package common

import (
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logTimeFormat = "2006-01-02 15:04:05"
)

func InitLog() {
	core := zapcore.NewTee(
		zapcore.NewCore(logEncoder(), logFileWriteSyncer(), zapcore.InfoLevel),
		zapcore.NewCore(logEncoder(), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zapcore.DebugLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(logger)
}

func logEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    logEncodeLevel,
			EncodeTime:     logEncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   logEncodeCaller,
		})
}

func logFileWriteSyncer() zapcore.WriteSyncer {
	logPath := viper.GetString("log.logPath")
	maxSize := viper.GetInt("log.maxFileSize")
	maxBackups := viper.GetInt("log.maxBackups")
	maxAge := viper.GetInt("log.maxAge")
	compress := viper.GetBool("log.compress")

	writeSyncer := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(
			&lumberjack.Logger{
				Filename:   logPath,
				MaxSize:    maxSize,
				MaxBackups: maxBackups,
				MaxAge:     maxAge,
				Compress:   compress,
			}),
		Size: 4096,
	}

	return zapcore.AddSync(writeSyncer)
}

func logEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")

}

func logEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(logTimeFormat) + "]")
}

func logEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}
