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
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"oset/model"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type Claims struct {
	Uid    int
	Level  model.UserLevel
	Uname  string
	Email  string
	Avatar string
	jwt.RegisteredClaims
}

var (
	jwtKey     []byte
	jwtKeyInit sync.Once
)

func JwtKey() []byte {
	jwtKeyInit.Do(func() {
		jwtKey = []byte(viper.GetString("sys.jwt_key"))
		if len(jwtKey) == 0 {
			rand.Seed(time.Now().UnixNano())
			sbytes := make([]byte, 256)
			rand.Read(sbytes)

			hash := sha256.New()
			hash.Write(sbytes)
			jwtKeyS := hex.EncodeToString(hash.Sum(nil))
			jwtKey = []byte(jwtKeyS)
			viper.Set("sys.jwt_key", jwtKeyS)
		}
	})

	return jwtKey
}

func GenerateToken(user *model.User) (token string, err error) {
	nowTime := time.Now()
	// token 有效期30天
	expireTime := nowTime.Add(30 * 24 * time.Hour)
	claims := &Claims{
		Uid:    user.Uid,
		Level:  user.Level,
		Uname:  user.Uname,
		Email:  user.Email,
		Avatar: user.Avatar,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "oset",
		},
	}

	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = _token.SignedString(JwtKey())
	return
}

func ParseToken(tokenString string) (token *jwt.Token, claims *Claims, err error) {
	claims = &Claims{}
	token, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return JwtKey(), nil
	})

	return
}
