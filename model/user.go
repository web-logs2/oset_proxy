//
// File: user.go
// Created by Dizzrt on 2023/01/17.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package model

type RoleType int

const (
	ROLETYPE_NORMAL = iota
	ROLETYPE_ADMIN
)

type User struct {
	Uid       int      `gorm:"primaryKey" json:"uid" form:"uid"`
	Role      RoleType `gorm:"not null" json:"role" form:"role"`
	Uname     string   `gorm:"size:32;not null" json:"uname" form:"uname"`
	Password  string   `gorm:"size:255;not null" json:"password" form:"password"`
	Email     string   `gorm:"size:64;not null" json:"email" form:"email"`
	Avatar    string   `gorm:"size:255;not null" json:"avatar" form:"avatar"`
	Activated bool     `gorm:"bool;default:false" json:"activated" form:"activated"`
	CreatedAt int
	UpdatedAt int
}

type UserInfo struct {
	Uid       int      `gorm:"primaryKey" json:"uid" form:"uid"`
	Role      RoleType `gorm:"not null" json:"role" form:"role"`
	Uname     string   `gorm:"size:32;not null" json:"uname" form:"uname"`
	Email     string   `gorm:"size:64;not null" json:"email" form:"email"`
	Avatar    string   `gorm:"size:255;not null" json:"avatar" form:"avatar"`
	Activated bool     `gorm:"bool;default:false" json:"activated" form:"activated"`
}
