//
// File: life.go
// Created by Dizzrt on 2023/02/28.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package main

import (
	"fmt"
	"oset/api/controller"
	"oset/component/log"
	"oset/db"
	"time"

	"github.com/Dizzrt/etfoundation/sig"
	"github.com/Dizzrt/etlog"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName(".config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic("read config failed: " + err.Error())
	}

	log.InitLog()
	controller.InitEvent()
	db.InitMysqlFromViper()
	db.InitRedisFromViper()
}

func Defer() {
	fmt.Printf("\nstopping oset...\n")

	etlog.L().Sync()
	viper.WriteConfig()
	time.Sleep(time.Second)

	sig.DoQuit()
	time.Sleep(time.Second)

	fmt.Println("successfully stopped oset")
}
