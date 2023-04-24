//
// File: auth.go
// Created by Dizzrt on 2023/04/21.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

//
// File: auth.go
// Created by Dizzrt on 2023/04/21.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package model

import "encoding/json"

type AKSK struct {
	Ak         string `gorm:"char(64);not null" json:"ak"`
	Sk         string `gorm:"char(64);not null" json:"sk"`
	Aid        int    `gorm:"index;not null" json:"aid"`
	ExpireTime int64  `gorm:"default:0" json:"expire_time"`
}

func (aksk AKSK) MarshalBinary() ([]byte, error) {
	return json.Marshal(aksk)
}

func (aksk *AKSK) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, aksk)
}

type AKSKExtension struct {
	ID int `gorm:"primaryKey" json:"id"`
	AKSK
	Description string
	CreatedAt   int
	UpdatedAt   int
}
