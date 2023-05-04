//
// File: user.go
// Created by Dizzrt on 2023/03/01.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package controller

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"oset/auth"
	"oset/common"
	"oset/db"
	"oset/model"
	"regexp"
	"strconv"
	"time"

	"github.com/Dizzrt/etlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func abortCtx(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"msg": msg,
	})
	ctx.Abort()
}

func Login(ctx *gin.Context) {
	var requestUser model.User
	ctx.BindJSON(&requestUser)

	var user model.User
	res := db.Mysql().First(&user, "email = ?", requestUser.Email)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		abortCtx(ctx, http.StatusOK, "用户不存在")
		return
	}

	hashedPassword := []byte(user.Password)
	isPwdMatched := bcrypt.CompareHashAndPassword(hashedPassword, []byte(requestUser.Password))

	if isPwdMatched != nil {
		if errors.Is(isPwdMatched, bcrypt.ErrMismatchedHashAndPassword) {
			abortCtx(ctx, http.StatusOK, "密码错误")
		} else {
			abortCtx(ctx, http.StatusInternalServerError, isPwdMatched.Error())
		}

		return
	}

	token, err := auth.GenerateToken(&user)
	if err != nil {
		abortCtx(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "登陆成功",
		"token":  token,
		"name":   user.Uname,
		"email":  user.Email,
		"uid":    user.Uid,
		"avatar": user.Avatar,
	})
}

func CreateUser(ctx *gin.Context) {
	ru, _ := ctx.Get("user")
	requestUser := ru.(model.User)

	if requestUser.Level < model.USERLEVEL_ADMIN {
		abortCtx(ctx, http.StatusUnauthorized, "权限不足")
		return
	}

	var newUser model.User
	ctx.BindJSON(&newUser)

	res := db.Mysql().First(&newUser, "email = ?", newUser.Email)
	if res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			etlog.L().Error(res.Error.Error(), zap.Int("request_user_id", requestUser.Uid))
			abortCtx(ctx, http.StatusInternalServerError, "unknown error")
			return
		}
	}

	if res.Error == nil {
		abortCtx(ctx, http.StatusOK, "邮箱已被占用")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		etlog.L().Error(err.Error(), zap.Int("request_user_id", requestUser.Uid))
		abortCtx(ctx, http.StatusInternalServerError, "unknown error")
		return
	}

	var uid int
	rand.Seed(time.Now().Unix())
	for {
		uid = rand.Intn(8999999) + 1000000
		res = db.Mysql().First(&newUser, uid)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			break
		}
	}

	newUser.Uid = uid
	newUser.Password = string(hashedPassword)
	if newUser.Avatar == "" {
		newUser.Avatar = "/images/avatar.png"
	}

	res = db.Mysql().Create(&newUser)
	if res.Error != nil {
		abortCtx(ctx, http.StatusInternalServerError, "unknown error")
		return
	}

	jsonBytes, err := json.Marshal(newUser)
	if err != nil {
		abortCtx(ctx, http.StatusInternalServerError, "unknown error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "创建成功",
		"user_info": string(jsonBytes),
	})
}

func GetUserList(ctx *gin.Context) {
	ru, _ := ctx.Get("user")
	requestUser := ru.(model.User)

	if requestUser.Level < model.USERLEVEL_ADMIN {
		abortCtx(ctx, http.StatusUnauthorized, "权限不足")
		return
	}

	var userList []model.UserInfo
	res := db.Mysql().Table("users").Find(&userList)

	if res.Error != nil {
		etlog.L().Error("failed to get user list", zap.Error(res.Error))
		abortCtx(ctx, http.StatusInternalServerError, "unknown error")
		return
	}

	jsonBytes, err := json.Marshal(userList)
	if err != nil {
		etlog.L().Error(err.Error())
		abortCtx(ctx, http.StatusInternalServerError, "unknown error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "success",
		"user_list": string(jsonBytes),
	})
}

func GetUserInfo(ctx *gin.Context) {
	ru, _ := ctx.Get("user")
	requestUser := ru.(model.User)

	target, isExist := ctx.GetQuery("target")
	if !isExist {
		abortCtx(ctx, http.StatusOK, "no target")
		return
	}

	pattern := `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	match, _ := regexp.MatchString(pattern, target)

	var res *gorm.DB
	var targetUser model.UserInfo

	if match {
		// email
		res = db.Mysql().Table("users").Where("email = ?", target).First(&targetUser)
	} else if _, err := strconv.Atoi(target); err == nil {
		// uid
		res = db.Mysql().Table("users").Where("uid = ?", target).First(&targetUser)
	} else {
		// uname
		res = db.Mysql().Table("users").Where("uname = ?", target).First(&targetUser)
	}

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"msg":  "the user does not exist",
				"data": "{}",
			})
			return
		} else {
			etlog.L().Error(res.Error.Error())
			abortCtx(ctx, http.StatusInternalServerError, "unknown error")
			return
		}
	}

	if requestUser.Level < model.USERLEVEL_ADMIN && requestUser.Uid != targetUser.Uid {
		abortCtx(ctx, http.StatusUnauthorized, "权限不足")
		return
	}

	jsonBytes, err := json.Marshal(targetUser)
	if err != nil {
		abortCtx(ctx, http.StatusInternalServerError, "unknown error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": string(jsonBytes),
	})
}

func SetUserInfo(ctx *gin.Context) {
	ru, _ := ctx.Get("user")
	requestUser := ru.(model.User)

	var targetUser model.UserInfo
	ctx.BindJSON(&targetUser)

	if requestUser.Level < model.USERLEVEL_ADMIN && requestUser.Uid != targetUser.Uid {
		abortCtx(ctx, http.StatusUnauthorized, "权限不足")
		return
	}

	res := db.Mysql().Model(&model.User{}).Where("uid = ?", targetUser.Uid).Updates(map[string]interface{}{
		"uname":     targetUser.Uname,
		"email":     targetUser.Email,
		"activated": targetUser.Activated,
	})

	if res.Error != nil {
		etlog.L().Error(res.Error.Error())
		abortCtx(ctx, http.StatusInternalServerError, "unknown error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func DropUser(ctx *gin.Context) {
	ru, _ := ctx.Get("user")
	requestUser := ru.(model.User)

	if requestUser.Level < model.USERLEVEL_ADMIN {
		abortCtx(ctx, http.StatusUnauthorized, "权限不足")
		return
	}

	targetUid, err := strconv.Atoi(ctx.Query("uid"))
	if err != nil {
		etlog.L().Error(err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusCommonError,
			"msg":  "error",
		})
		ctx.Abort()
		return
	}

	res := db.Mysql().Delete(&model.User{}, targetUid)
	if res.Error != nil {
		etlog.L().Error(res.Error.Error())
		abortCtx(ctx, http.StatusInternalServerError, "unknown error")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
