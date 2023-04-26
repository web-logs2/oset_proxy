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
	"os"
	"oset/common"

	"github.com/Dizzrt/etlog"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

func UploadImg(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})

		ctx.Abort()
		etlog.L().Error("upload img failed", zap.Error(err))
		return
	}

	if _, err := os.Stat("./static/upload/image"); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll("./static/upload/image", 0777)
		}
	}

	uuid := uuid.New()
	fileName := uuid.String() + ".png"
	if err := ctx.SaveUploadedFile(file, "./static/upload/image/"+fileName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})

		ctx.Abort()
		etlog.L().Error("save img failed", zap.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"img": viper.GetString("sys.self_host") + "/static/stream/" + fileName,
	})
}

func GetUploadImgStream(ctx *gin.Context) {
	img := ctx.Param("img")
	file, _ := os.ReadFile("./static/upload/image/" + img)
	ctx.Writer.WriteString(string(file))
}
