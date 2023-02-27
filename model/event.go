//
// File: event.go
// Created by Dizzrt on 2023/02/24.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package model

import "time"

type Event struct {
	Aid   int       `json:"aid" form:"aid"`
	Did   int       `json:"did" form:"did"`
	Event string    `json:"event" form:"event"`
	Data  string    `json:"data" form:"data"`
	Time  time.Time `json:"time" form:"time"`
}
