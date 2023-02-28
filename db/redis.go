//
// File: redis.go
// Created by Dizzrt on 2023/02/28.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package db

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	redisDB *redis.Client
)

func InitRedis(addr string, password string, db int) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	redisDB = rdb
}

func InitRedisFromViper() {
	addr := viper.GetString("redis.host") + ":" + viper.GetString("redis.port")
	pwd := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")

	InitRedis(addr, pwd, db)
}

func Redis() *redis.Client {
	return redisDB
}
