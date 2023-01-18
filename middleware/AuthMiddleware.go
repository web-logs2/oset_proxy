//
// File: AuthMiddleware.go
// Created by Dizzrt on 2023/01/18.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"oset/common"
	"oset/common/auth"
	"oset/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")

		if tokenString == "" || len(tokenString) < 7 || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": common.StatusTokenMalformed,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := auth.ParseToken(tokenString)

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

		user := model.User{
			Uid:    claims.Uid,
			Uname:  claims.Uname,
			Email:  claims.Email,
			Avatar: claims.Avatar,
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}
