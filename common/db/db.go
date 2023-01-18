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

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gDB *gorm.DB

func InitDB(host string, port string, user string, password string, db string, charset string, loc string) *gorm.DB {
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true&loc=%s", user, password, host, port, db, charset, url.QueryEscape(loc))
	newDB, err := gorm.Open(mysql.Open(args), &gorm.Config{})

	if err != nil {
		panic("failed to open database: " + err.Error())
	}

	gDB = newDB
	DBMigrate()

	return gDB
}

func DBMigrate() {
	gDB.AutoMigrate(&model.User{})
}

func GetDB() *gorm.DB { return gDB }
