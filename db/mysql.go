//
// File: mysql.go
// Created by Dizzrt on 2023/02/28.
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

var (
	mysqlDB *gorm.DB
)

func InitMysql(host string, port string, database string, user string, password string, loc string, charset string) {
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true&loc=%s", user, password, host, port, database, charset, url.QueryEscape(loc))

	mdb, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		etlog.L().Panic("failed to connect to mysql", zap.Error(err))
	}

	mysqlDB = mdb
	AutoMigrate()
}

func InitMysqlFromViper() {
	mysqlCfg := viper.GetStringMapString("mysql")

	host := mysqlCfg["host"]
	port := mysqlCfg["port"]
	database := mysqlCfg["database"]
	user := mysqlCfg["user"]
	pwd := mysqlCfg["password"]
	loc := mysqlCfg["loc"]
	charset := mysqlCfg["charset"]

	InitMysql(host, port, database, user, pwd, loc, charset)
}

func AutoMigrate() {
	err := mysqlDB.AutoMigrate(&model.User{})
	if err != nil {
		etlog.L().Panic("failed to migrate user table", zap.Error(err))
	}

	err = mysqlDB.Set("gorm:table_options", "AUTO_INCREMENT=1001").AutoMigrate(&model.App{})
	if err != nil {
		etlog.L().Panic("failed to migrate app table", zap.Error(err))
	}
}

func Mysql() *gorm.DB {
	return mysqlDB
}
