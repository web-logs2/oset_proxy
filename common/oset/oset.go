//
// File: init.go
// Created by Dizzrt on 2023/01/18.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package oset

import (
	"fmt"
	"oset/common/component/log"
	"oset/common/db"
	"oset/common/oset/info"
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

	info.Init()
	// stream.InitKafka()
	log.InitLog()
	db.InitDB()
}

func Defer() {
	fmt.Printf("\nstopping oset...\n")

	zap.L().Sync()
	time.Sleep(time.Second)

	close(info.QuitSig)
	time.Sleep(time.Second)

	fmt.Println("successfully stopped oset")
}
