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
	"oset/api"
	"oset/api/controller"
	"oset/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoutes(r *gin.Engine) *gin.Engine {
	r.Use(middleware.GinLogger())
	r.Use(middleware.CORSMiddleware())

	r.Static("/static/upload/image", "./static/upload/image")
	r.POST("/static/upload/image", api.UploadImg)
	r.GET("/static/stream/:img", api.GetUploadImgStream)

	sysRoutes := r.Group("/sys")
	sysRoutes.GET("/getinit", api.GetInit)
	sysRoutes.POST("/setinit", api.SetInit)

	r.POST("/login", controller.Login)

	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.JwtMiddleware())
	userRoutes.GET("list", controller.GetUserList)
	userRoutes.GET("info", controller.GetUserInfo)
	userRoutes.POST("update", controller.SetUserInfo)
	userRoutes.POST("create", controller.CreateUser)
	userRoutes.DELETE("drop", controller.DropUser)

	appRoutes := r.Group("/app")
	appRoutes.Use(middleware.JwtMiddleware())
	appRoutes.GET("info", controller.GetApp)
	appRoutes.POST("create", controller.CreateApp)
	appRoutes.POST("update", controller.UpdateApp)
	appRoutes.DELETE("delete", controller.DropApp)
	appRoutes.GET("list", controller.GetAppList)
	appRoutes.GET("aksk/list", controller.GetAppAkSK)
	appRoutes.POST("aksk/generate", controller.GenerateAKSK)
	appRoutes.POST("aksk/update", controller.UpdateAksk)
	appRoutes.DELETE("aksk/delete", controller.DropAKSK)

	eventRoutes := r.Group("/event")
	eventRoutes.POST("report/:aid", middleware.AkskMiddleware(), controller.ReportEvent)
	eventRoutes.GET("tool/realtime/:aid/:did", controller.RegisterRealtimeEvent)
	return r
}
