//
// File: app.go
// Created by Dizzrt on 2023/02/01.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package model

type App struct {
	Aid         int    `gorm:"primaryKey;autoIncrement" json:"aid" form:"aid"`
	Icon        string `gorm:"size:255;not null" json:"icon" form:"icon"`
	Name        string `gorm:"size:32;not null" json:"name" form:"name"`
	Description string `gorm:"size:255;" json:"des" form:"des"`
	Activated   bool   `gorm:"bool;default:false" json:"activated" form:"activated"`
	CreatedAt   int
	UpdatedAt   int
}
