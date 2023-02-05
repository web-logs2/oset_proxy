//
// File: net_stream.go
// Created by Dizzrt on 2023/02/05.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package stream

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

func GetRawBody(ctx *gin.Context) (bodyBytes []byte, err error) {
	bodyBytes, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return
}
