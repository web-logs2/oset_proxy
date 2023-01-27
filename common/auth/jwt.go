//
// File: jwt.go
// Created by Dizzrt on 2023/01/18.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package auth

import (
	"oset/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("4420335dd9486c4c41f8d21176e5d37399502d14")

type Claims struct {
	Uid    int
	Role   model.RoleType
	Uname  string
	Email  string
	Avatar string
	jwt.RegisteredClaims
}

func GenerateToken(user *model.User) (token string, err error) {
	nowTime := time.Now()
	// token 有效期30天
	expireTime := nowTime.Add(30 * 24 * time.Hour)
	claims := &Claims{
		Uid:    user.Uid,
		Role:   user.Role,
		Uname:  user.Uname,
		Email:  user.Email,
		Avatar: user.Avatar,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "oset",
		},
	}

	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = _token.SignedString(jwtKey)
	return
}

func ParseToken(tokenString string) (token *jwt.Token, claims *Claims, err error) {
	claims = &Claims{}
	token, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return
}
