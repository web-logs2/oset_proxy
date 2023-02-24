//
// File: log.go
// Created by Dizzrt on 2023/02/24.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package log

import (
	"sync"

	"github.com/Dizzrt/etlog"
	"github.com/Dizzrt/etstream/kafka"
	"github.com/spf13/viper"
)

var (
	once sync.Once
)

func InitLog() {
	once.Do(func() {
		reporterType := viper.GetString("log.reporter_type")
		reportername := viper.GetString("log.reporter_name")
		logFilePath := viper.GetString("log.file_path")
		maxFileSize := viper.GetInt("log.max_file_size")
		maxBackups := viper.GetInt("log.max_backups")
		maxAge := viper.GetInt("log.max_age")
		compress := viper.GetBool("log.is_compress")

		kafkaEnable := viper.GetBool("log.kafka.is_enable")
		kafkaHost := viper.GetString("log.kafka.host")
		kafkaTopic := viper.GetString("log.kafka.topic")

		kconfig := kafka.KafkaConfig{
			SaramaConfig: kafka.DefaultProducerConfig(),
			Host:         kafkaHost,
			Topic:        kafkaTopic,
		}

		lconfig := etlog.LogConfig{
			ReporterType: reporterType,
			ReporterName: reportername,
			FilePath:     logFilePath,
			MaxFileSize:  maxFileSize,
			MaxBackups:   maxBackups,
			MaxAge:       maxAge,
			Compress:     compress,
			KafkaEnable:  kafkaEnable,
			KafkaConfig:  kconfig,
		}

		etlog.NewLogger(lconfig, "proxy_log")
	})
}
