//
// File: user_controller.go
// Created by Dizzrt on 2023/01/17.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package controller

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"oset/common"
	"oset/common/auth"
	"oset/common/db"
	"oset/model"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	_db := db.GetDB()

	var requestUser model.User
	ctx.Bind(&requestUser)

	var user model.User
	resDB := _db.First(&user, "email = ?", requestUser.Email)
	if !errors.Is(resDB.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": common.StatusEmailUsed,
			"msg":  "邮箱已被占用",
		})
		ctx.Abort()
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(requestUser.Password), bcrypt.DefaultCost)

	var uid int
	user.Uid = 0
	rand.Seed(time.Now().Unix())
	for {
		// uid range 1000000~9999999
		uid = rand.Intn(8999999) + 1000000
		resDB = _db.First(&user, uid)
		if errors.Is(resDB.Error, gorm.ErrRecordNotFound) {
			break
		}
	}

	user = model.User{
		Uid:      uid,
		Role:     requestUser.Role,
		Uname:    requestUser.Uname,
		Password: string(hashedPassword),
		Email:    requestUser.Email,
		Avatar:   "/images/avatar.png",
	}

	_db.Create(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusRegisterOk,
		"msg":  "注册成功",
	})
}

func Login(ctx *gin.Context) {
	_db := db.GetDB()

	var requestUser model.User
	ctx.Bind(&requestUser)

	var user model.User
	resDB := _db.First(&user, "email = ?", requestUser.Email)
	if errors.Is(resDB.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": common.StatusUserNotExist,
			"msg":  "用户不存在",
		})
		ctx.Abort()
		return
	}

	hashedPassword := []byte(user.Password)
	isPwdMatched := bcrypt.CompareHashAndPassword(hashedPassword, []byte(requestUser.Password))

	if isPwdMatched != nil {
		if errors.Is(isPwdMatched, bcrypt.ErrMismatchedHashAndPassword) {
			ctx.JSON(http.StatusOK, gin.H{
				"code": common.StatusWrongPassword,
				"msg":  "密码错误",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"code": common.StatusUserUnhandled,
				"msg":  isPwdMatched.Error(),
			})
		}

		ctx.Abort()
		return
	}

	token, err := auth.GenerateToken(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusGenTokenError,
			"msg":  fmt.Sprintf("[%d]server error: %s", common.StatusGenTokenError, err.Error()),
		})

		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusLoginOk,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})
}
