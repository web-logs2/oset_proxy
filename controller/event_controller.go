//
// File: et_controller.go
// Created by Dizzrt on 2023/01/31.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package controller

import (
	"io"
	"net/http"
	"oset/common/stream"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ReportEvent(ctx *gin.Context) {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("failed to get request body", zap.String("err", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "error",
		})
		ctx.Abort()
		return
	}

	stream.SendEvent(string(bodyBytes))

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
