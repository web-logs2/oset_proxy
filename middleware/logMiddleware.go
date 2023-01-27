//
// File: logMiddleware.go
// Created by Dizzrt on 2023/01/25.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package middleware

import (
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		bodyBytes, _ := ioutil.ReadAll(ctx.Request.Body)
		body := string(bodyBytes)
		start := time.Now()
		ctx.Next()
		elapsed := time.Since(start).Milliseconds()
		zap.L().Info(path,
			zap.Int("status", ctx.Writer.Status()),
			zap.String("method", ctx.Request.Method),
			zap.String("query", query),
			zap.String("body", body),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Int64("elapsed", elapsed),
		)
	}
}
