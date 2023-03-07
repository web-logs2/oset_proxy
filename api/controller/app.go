//
// File: app.go
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
	"oset/db"

	"oset/model"
	"strconv"

	"github.com/Dizzrt/etlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// func authVerify(role model.RoleType) bool {

// 	return true
// }

func CreateApp(ctx *gin.Context) {
	requestUser, _ := ctx.Get("user")
	ruser := requestUser.(model.User)

	newApp := model.App{}
	err := ctx.BindJSON(&newApp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "create new app failed",
		})
		ctx.Abort()

		etlog.L().Error("unable to create new app, because bindjson failed", zap.Int("operator_uid", ruser.Uid), zap.Error(err))
		return
	}

	if newApp.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "create new app failed, the new app has no name",
		})
		ctx.Abort()

		etlog.L().Warn("unable to create new app, because app name is empty", zap.Any("new_app", newApp), zap.Int("operatpr_uid", ruser.Uid))
		return
	}

	if newApp.Icon == "" {
		newApp.Icon = "/images/defaultIcon.png"
	}

	res := db.Mysql().Create(&newApp)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "failed to create App, " + res.Error.Error(),
		})
		ctx.Abort()

		etlog.L().Error("unable to create new app", zap.Error(res.Error))
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": common.StatusCommonOK,
			"msg":  newApp.Name + " successfully created",
		})
		etlog.L().Info("created new app", zap.String("app_name", newApp.Name), zap.Int("operator_uid", ruser.Uid))
	}
}

func UpdateApp(ctx *gin.Context) {
	requestUser, _ := ctx.Get("user")
	ruser := requestUser.(model.User)

	var newAppInfo model.App
	ctx.BindJSON(&newAppInfo)

	res := db.Mysql().Model(&model.App{}).Where("aid = ?", newAppInfo.Aid).Updates(model.App{
		Icon:        newAppInfo.Icon,
		Name:        newAppInfo.Name,
		Description: newAppInfo.Description,
	})

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "failed to update App info",
		})

		etlog.L().Error("update App info failed", zap.Int("aid", newAppInfo.Aid), zap.Any("new_info", newAppInfo), zap.Int("uid", ruser.Uid), zap.Error(res.Error))
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
		etlog.L().Error("delete app failed", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("err", err.Error()))

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "error",
		})
		ctx.Abort()
		return
	}

	res := db.Mysql().Delete(&model.App{}, aid)
	if res.Error != nil {
		etlog.L().Error("delete app failed", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("err", res.Error.Error()))

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
		etlog.L().Error("get app info failed", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("err", err.Error()))

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "error",
		})
		ctx.Abort()
		return
	}

	var targetApp model.App
	res := db.Mysql().Where("aid = ?", aid).First(&targetApp)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"code": common.StatusCommonOK,
				"data": "{}",
				"msg":  "the app does not exist",
			})
			return
		} else {
			etlog.L().Error("search app error", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("err", res.Error.Error()))

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

		etlog.L().Error("search app error", zap.Int("aid", aid), zap.Int("uid", requestUser.(model.User).Uid), zap.String("error", res.Error.Error()))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusCommonOK,
		"data": string(jsonBytes),
		"msg":  "success",
	})
}
