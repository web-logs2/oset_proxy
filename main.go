// File: main.go
// Created by Dizzrt on 2023/01/17.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"oset/common/db"
	"oset/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// load config
	configFile, err := os.Open("./config.json")
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	configMap := make(map[string]interface{})
	err = json.NewDecoder(configFile).Decode(&configMap)
	if err != nil {
		panic("failed to decode config: " + err.Error())
	}
	configFile.Close()

	// init database
	{
		host := configMap["db_host"].(string)
		port := configMap["db_port"].(string)
		user := configMap["db_user"].(string)
		pwd := configMap["db_password"].(string)
		dbName := configMap["db_name"].(string)
		charset := configMap["db_charset"].(string)
		loc := configMap["db_loc"].(string)
		db.InitDB(host, port, user, pwd, dbName, charset, loc)
	}

	// init router
	r := gin.Default()
	r.StaticFS("/image", http.Dir("./static/image"))
	router.CollectRoutes(r)
	panic(r.Run())
}
