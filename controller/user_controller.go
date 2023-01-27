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
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"oset/common"
	"oset/common/auth"
	"oset/common/db"
	"oset/model"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func GetUserList(ctx *gin.Context) {
	requestUser, _ := ctx.Get("user")
	if requestUser.(model.User).Role < model.ROLETYPE_ADMIN {
		zap.L().Warn("trying to get user list, but permission denied.", zap.Int("uid", requestUser.(model.User).Uid), zap.String("user", requestUser.(model.User).Uname))

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": common.StatusUserLowPermission,
			"data": "{}",
			"msg":  "error",
		})
		ctx.Abort()
		return
	}

	var users []model.UserInfo
	resDB := db.GetDB().Table("users").Find(&users)

	if resDB.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusUserError,
			"data": "{}",
			"msg":  "error",
		})

		zap.L().Error("failed to query all user list.", zap.String("error", resDB.Error.Error()))
		ctx.Abort()
		return
	}

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusUserError,
			"data": "{}",
			"msg":  "error",
		})

		zap.L().Error("failed to query all user list.", zap.String("error", resDB.Error.Error()))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusUserOk,
		"data": string(jsonBytes),
		"msg":  "success",
	})
}

func GetUserInfo(ctx *gin.Context) {
	requestUser, _ := ctx.Get("user")
	target, isExist := ctx.GetQuery("target")

	if !isExist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": common.StatusUserError,
			"data": "{}",
			"msg":  "no target",
		})
		ctx.Abort()
		return
	}

	if requestUser.(model.User).Role < model.ROLETYPE_ADMIN {
		zap.L().Warn("trying to get user info, but permission denied.", zap.Int("uid", requestUser.(model.User).Uid), zap.String("user", requestUser.(model.User).Uname), zap.String("targetUser", target))

		ctx.JSON(http.StatusOK, gin.H{
			"code": common.StatusUserLowPermission,
			"data": "{}",
			"msg":  "error",
		})
		ctx.Abort()
		return
	}

	pattern := `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	match, _ := regexp.MatchString(pattern, target)

	var resDB *gorm.DB
	var targetUser model.UserInfo

	if match {
		// email
		resDB = db.GetDB().Table("users").Where("email = ?", target).First(&targetUser)
	} else if _, err := strconv.Atoi(target); err == nil {
		// uid
		resDB = db.GetDB().Table("users").Where("uid = ?", target).First(&targetUser)
	} else {
		// uname
		resDB = db.GetDB().Table("users").Where("uname = ?", target).First(&targetUser)
	}

	if resDB.Error != nil {
		if errors.Is(resDB.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"code": common.StatusUserOk,
				"data": "{}",
				"msg":  "the user does not exist",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": common.StatusUserError,
				"data": "{}",
				"msg":  "error",
			})

			zap.L().Error("search user error", zap.String("target", target), zap.String("error", resDB.Error.Error()))
			ctx.Abort()
			return
		}
	}

	jsonBytes, err := json.Marshal(targetUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusUserError,
			"data": "{}",
			"msg":  "error",
		})

		zap.L().Error("search user error", zap.String("target", target), zap.String("error", resDB.Error.Error()))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusUserOk,
		"data": string(jsonBytes),
		"msg":  "success",
	})
}

func SetUserInfo(ctx *gin.Context) {
	requestUser, _ := ctx.Get("user")
	var targetUser model.UserInfo
	ctx.Bind(&targetUser)

	if requestUser.(model.User).Role < model.ROLETYPE_ADMIN && targetUser.Uid != requestUser.(model.User).Uid {
		jsonBytes, _ := json.Marshal(targetUser)
		zap.L().Warn("trying to update user info, but permission denied.", zap.Int("uid", requestUser.(model.User).Uid), zap.String("user", requestUser.(model.User).Uname), zap.Int("target_user", targetUser.Uid), zap.String("target_value", string(jsonBytes)))

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": common.StatusUserError,
			"msg":  "error",
		})

		ctx.Abort()
		return
	}

	resDB := db.GetDB().Model(&model.User{}).Where("uid = ?", targetUser.Uid).Updates(model.UserInfo{
		Uname:  targetUser.Uname,
		Email:  targetUser.Email,
		Avatar: targetUser.Avatar,
	})

	if resDB.Error != nil {
		zap.L().Error("update user info failed", zap.Int("target_user", targetUser.Uid), zap.String("target_uname", targetUser.Uname), zap.String("target_email", targetUser.Email), zap.String("traget_avatar", targetUser.Avatar), zap.String("error", resDB.Error.Error()))

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusUserError,
			"msg":  "error",
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": common.StatusUserOk,
		"msg":  "success",
	})
}
