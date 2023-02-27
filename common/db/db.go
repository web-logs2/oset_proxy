//
// File: db.go
// Created by Dizzrt on 2023/01/17.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package db

import (
	"fmt"
	"net/url"
	"oset/model"

	"github.com/Dizzrt/etlog"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gDB *gorm.DB

func InitDB() {
	db_config := viper.GetStringMapString("mysql")
	user := db_config["user"]
	password := db_config["password"]
	host := db_config["host"]
	port := db_config["port"]
	db := db_config["database"]
	charset := db_config["charset"]
	loc := db_config["loc"]

	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true&loc=%s", user, password, host, port, db, charset, url.QueryEscape(loc))
	newDB, err := gorm.Open(mysql.Open(args), &gorm.Config{})

	if err != nil {
		etlog.L().Panic("failed to connect to database", zap.String("error", err.Error()))
	}

	gDB = newDB
	DBMigrate()
}

func DBMigrate() {
	err := gDB.AutoMigrate(&model.User{})
	if err != nil {
		etlog.L().Panic("failed to migrate user table", zap.Error(err))
	}

	err = gDB.Set("gorm:table_options", "AUTO_INCREMENT=1001").AutoMigrate(&model.App{})
	if err != nil {
		etlog.L().Panic("failed to migrate app table", zap.Error(err))
	}
}

func GetDB() *gorm.DB { return gDB }
