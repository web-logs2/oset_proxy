//
// File: app_controller.go
// Created by Dizzrt on 2023/02/01.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"oset/common"
	"oset/common/db"
	"oset/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// func authVerify(role model.RoleType) bool {

// 	return true
// }

func CreateApp(ctx *gin.Context) {
	requestUser, _ := ctx.Get("user")

	var newApp model.App
	ctx.Bind(&newApp)

	_db := db.GetDB()
	res := _db.Create(&newApp)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "failed to create App, " + res.Error.Error(),
		})
		zap.L().Error("failed to create App", zap.String("err", res.Error.Error()))
		ctx.Abort()
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": common.StatusCommonOK,
			"msg":  newApp.Name + " successfully created",
		})
		zap.L().Info("create new app", zap.String("app_name", newApp.Name), zap.Int("creator", requestUser.(model.User).Uid))
	}
}

func UpdateApp(ctx *gin.Context) {
	requestUser, _ := ctx.Get("user")

	var newAppInfo model.App
	ctx.Bind(&newAppInfo)

	res := db.GetDB().Model(&model.App{}).Where("aid = ?", newAppInfo.Aid).Updates(model.App{
		Icon:        newAppInfo.Icon,
		Name:        newAppInfo.Name,
		Description: newAppInfo.Description,
	})

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "failed to update App info",
		})

		zap.L().Error("update App info failed", zap.Int("aid", newAppInfo.Aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("err", res.Error.Error()))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusCommonOK,
		"msg":  "success",
	})
}

func DropApp(ctx *gin.Context) {
	requestUser, _ := ctx.Get("user")
	aid, err := strconv.Atoi(ctx.Query("aid"))
	if err != nil {
		zap.L().Error("delete app failed", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("err", err.Error()))

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "error",
		})
		ctx.Abort()
		return
	}

	res := db.GetDB().Delete(&model.App{}, aid)
	if res.Error != nil {
		zap.L().Error("delete app failed", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("err", res.Error.Error()))

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "error",
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusCommonOK,
		"msg":  "success",
	})
}

func GetApp(ctx *gin.Context) {
	requestUser, _ := ctx.Get("user")
	aid, err := strconv.Atoi(ctx.Query("aid"))

	if err != nil {
		zap.L().Error("get app info failed", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("err", err.Error()))

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "error",
		})
		ctx.Abort()
		return
	}

	var targetApp model.App
	res := db.GetDB().Where("aid = ?", aid).First(&targetApp)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"code": common.StatusCommonOK,
				"data": "{}",
				"msg":  "the app does not exist",
			})
			return
		} else {
			zap.L().Error("search app error", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("err", res.Error.Error()))

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": common.StatusCommonError,
				"data": "{}",
				"msg":  "err",
			})
			ctx.Abort()
			return
		}
	}

	jsonBytes, err := json.Marshal(targetApp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusUserError,
			"data": "{}",
			"msg":  "error",
		})

		zap.L().Error("search app error", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("error", res.Error.Error()))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusCommonOK,
		"data": string(jsonBytes),
		"msg":  "success",
	})
}
