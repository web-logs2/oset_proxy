//
// File: sys.go
// Created by Dizzrt on 2023/01/21.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package api

import (
	"net/http"
	"oset/common"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/spf13/viper"
)

func GetInit(ctx *gin.Context) {
	var mp map[string]interface{} = make(map[string]interface{})

	mp["inited"] = viper.GetBool("sys.inited")

	data, err := json.Marshal(mp)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": common.StatusCommonError,
			"data": "{}",
			"msg":  "json.Marshal failed: " + err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusCommonOK,
		"data": string(data),
		"msg":  "successful",
	})
}

func SetInit(ctx *gin.Context) {
	mp := make(map[string]interface{})
	err := ctx.BindJSON(&mp)
	if err != nil {
		ctx.Abort()
		return
	}

	inited, ok := mp["inited"]
	if ok {
		viper.Set("sys.inited", inited)
	}

	viper.WriteConfig()
}
