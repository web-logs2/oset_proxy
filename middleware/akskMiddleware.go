//
// File: akskMiddleware.go
// Created by Dizzrt on 2023/03/09.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package middleware

import (
	"oset/auth"
	"strconv"
	"time"

	"github.com/Dizzrt/etlog"
	"github.com/gin-gonic/gin"
)

const (
	headerAccessKey = `x-auth-accesskey`
	headerSignature = `x-auth-signature`
	headerTimestamp = `x-auth-timestamp`
	headerContent   = `x-auth-content`
)

func AkskMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accesskey := ctx.GetHeader(headerAccessKey)
		signature := ctx.GetHeader(headerSignature)
		timestamp := ctx.GetHeader(headerTimestamp)
		content := ctx.GetHeader(headerContent)

		if accesskey == "" || timestamp == "" || signature == "" || content == "" {
			abortCtxWithUnauthorized(ctx)
			return
		}

		unixt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			abortCtxWithUnhandleError(ctx)
			etlog.L().Error(err.Error())
			return
		}

		t := time.Unix(unixt, 0)
		timeDelta := time.Since(t)
		if (timeDelta > 5*time.Minute) || (timeDelta < -1*time.Minute) {
			abortCtxWithUnauthorized(ctx)
			return
		}

		err = auth.ValidateSignature(accesskey, signature, content)
		if err != nil {
			abortCtxWithUnauthorized(ctx)
			return
		}

		ctx.Next()
	}
}
