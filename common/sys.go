//
// File: init.go
// Created by Dizzrt on 2023/01/18.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package common

import (
	"fmt"
	"oset/common/db"
	"oset/common/stream"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init() {
	viper.SetConfigName(".config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic("read config failed: " + err.Error())
	}

	stream.InitKafka()
	InitLog()
	db.InitDB()

	zap.L().Info("initalize oset successfully")
}

func Defer() {
	fmt.Println()

	zap.L().Sync()
	fmt.Println("waiting for sync logs to kafka in 1s")
	time.Sleep(1 * time.Second)

	fmt.Println("server exiting")
}
