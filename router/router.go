//
// File: router.go
// Created by Dizzrt on 2023/01/17.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package router

import (
	"net/http"
	"oset/controller"
	"oset/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoutes(r *gin.Engine) *gin.Engine {
	r.Use(middleware.GinLogger())
	r.Use(middleware.CORSMiddleware())
	r.StaticFS("/image", http.Dir("../static/image"))

	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	sysRoutes := r.Group("/sys")
	sysRoutes.GET("/getinit", controller.GetInit)
	sysRoutes.POST("/setinit", controller.SetInit)

	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware())
	userRoutes.GET("list", controller.GetUserList)
	userRoutes.GET("info", controller.GetUserInfo)
	userRoutes.POST("update", controller.SetUserInfo)

	appRoutes := r.Group("/app")
	appRoutes.Use(middleware.AuthMiddleware())
	appRoutes.GET("info", controller.GetApp)
	appRoutes.POST("create", controller.CreateApp)
	appRoutes.POST("update", controller.UpdateApp)
	appRoutes.DELETE("delete", controller.DropApp)

	return r
}
