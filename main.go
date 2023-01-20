// File: main.go
// Created by Dizzrt on 2023/01/17.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.

package main

import (
	"net/http"
	"oset/common"
	"oset/router"

	"github.com/gin-gonic/gin"
)

func main() {
	common.Init()

	// init router
	r := gin.Default()
	r.StaticFS("/image", http.Dir("./static/image"))
	router.CollectRoutes(r)
	panic(r.Run())
}
