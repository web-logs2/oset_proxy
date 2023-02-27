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
	"oset/controller"
	"time"

	"github.com/Dizzrt/etfoundation/sig"
	"github.com/Dizzrt/etlog"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName(".config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic("read config failed: " + err.Error())
	}

	controller.InitEvent()
	log.InitLog()
	db.InitDB()
}

func Defer() {
	fmt.Printf("\nstopping oset...\n")

	etlog.L().Sync()
	time.Sleep(time.Second)

	sig.DoQuit()
	time.Sleep(time.Second)

	fmt.Println("successfully stopped oset")
}
