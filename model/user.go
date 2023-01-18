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
	ROLETYPE_ADMIN = iota
	ROLETYPE_NORMAL
)

type User struct {
	Uid       int      `gorm:"primaryKey"`
	Role      RoleType `gorm:"not null"`
	Uname     string   `gorm:"varchar(32);not null"`
	Password  string   `gorm:"size:255;not null"`
	Email     string   `gorm:"varchar(64);not null"`
	Avatar    string   `gorm:"size:255;not null"`
	Activated bool     `gorm:"bool;default:false"`
	CreatedAt int
	UpdatedAt int
}
