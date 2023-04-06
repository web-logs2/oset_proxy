//
// File: jwtMiddleware.go
// Created by Dizzrt on 2023/01/18.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"oset/auth"
	"oset/common"
	"oset/db"
	"oset/model"
	"strconv"
	"strings"
	"time"

	"github.com/Dizzrt/etlog"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func abortCtx(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"msg": msg,
	})
	ctx.Abort()
}

func abortCtxWithUnauthorized(ctx *gin.Context) {
	abortCtx(ctx, http.StatusUnauthorized, "权限不足")
}

func abortCtxWithUnhandleError(ctx *gin.Context) {
	abortCtx(ctx, http.StatusInternalServerError, "authentication failed, unhandle error")
}

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")

		if tokenString == "" || len(tokenString) < 7 || !strings.HasPrefix(tokenString, "Bearer") {
			abortCtxWithUnauthorized(ctx)
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := auth.ParseToken(tokenString)

		// 检查token是否有效
		if !token.Valid {
			var code int
			if errors.Is(err, jwt.ErrTokenMalformed) {
				code = common.StatusTokenMalformed
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				code = common.StatusTokenExpired
			} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
				code = common.StatusTokenNotValidYet
			} else {
				code = common.StatusTokenUnhandled
			}

			var msg string
			if code == common.StatusTokenUnhandled {
				msg = fmt.Sprintf("[%d]权限不足: %s", common.StatusTokenUnhandled, err.Error())
			} else {
				msg = "权限不足"
			}

			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  msg,
			})
			ctx.Abort()
			return
		}

		// 检查是否处于启用状态
		isActive, err := checkIfActiveByUid(claims.Uid)
		if err != nil {
			etlog.L().Error(err.Error())
			abortCtxWithUnhandleError(ctx)
			return
		}

		if !isActive {
			abortCtxWithUnauthorized(ctx)
			return
		}

		user := model.User{
			Uid:    claims.Uid,
			Level:  claims.Level,
			Uname:  claims.Uname,
			Email:  claims.Email,
			Avatar: claims.Avatar,
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}

func checkIfActiveByUid(uid int) (isActive bool, err error) {
	rctx := context.Background()
	uidString := strconv.Itoa(uid)

	isActive, err = db.Redis().Get(rctx, uidString).Bool()
	// redis 中不存在，去 mysql 中查找并记录在 redis 中
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil

		res := db.Mysql().Model(&model.User{}).Select("activated").Where("uid = ?", uid).First(&isActive)
		if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			err = res.Error
			return
		}

		err = db.Redis().Set(rctx, uidString, isActive, time.Second*60).Err()
	}

	return
}
