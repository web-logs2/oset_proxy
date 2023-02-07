//
// File: info.go
// Created by Dizzrt on 2023/02/07.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package info

import (
	"github.com/spf13/viper"
)

const (
	ServerType = "proxy"
)

var (
	ServerName string
)

func Init() {
	ServerName = viper.GetString("sys.name")
}
