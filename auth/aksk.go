//
// File: aksk.go
// Created by Dizzrt on 2023/03/09.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"oset/db"
	"time"

	"github.com/Dizzrt/etlog"
	"gorm.io/gorm"
)

var (
	ErrSignatureInvalid  = errors.New("signature invalid")
	ErrNotFoundSecretKey = errors.New("secret not found")
	ErrAccessKeyExpired  = errors.New("access key has expired")
)

type AKSK struct {
	Ak         string `gorm:"primaryKey;char(64)" json:"ak"`
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
	AKSK
	Description string
	CreatedAt   int
	UpdatedAt   int
}

func GenerateAKSK(aid int, expireTime time.Duration, description string) (ak string, sk string, err error) {
	t := time.Now()

	var expireStamp int64
	if expireTime > 0 {
		expireStamp = t.Add(expireTime).Unix()
	} else {
		expireStamp = 0
	}

	tbytes := make([]byte, 8)
	binary.BigEndian.PutUint64(tbytes, uint64(t.Unix()))

	hash := sha256.New()
	hash.Write(tbytes)
	ak = hex.EncodeToString(hash.Sum(nil))

	hash.Write([]byte(ak))
	sk = hex.EncodeToString(hash.Sum(nil))

	aksk := AKSK{
		Ak:         ak,
		Sk:         sk,
		Aid:        aid,
		ExpireTime: expireStamp,
	}

	akskFull := AKSKExtension{
		AKSK:        aksk,
		Description: description,
	}

	res := db.Mysql().Create(&akskFull)
	if res.Error != nil {
		err = res.Error
		return
	}

	rctx := context.Background()
	err = db.Redis().Set(rctx, ak, aksk, 0).Err()
	if err != nil {
		return
	}
	err = db.Redis().ExpireAt(rctx, ak, time.Unix(expireStamp, 0)).Err()

	return
}

func getSK(ak string) (sk string, err error) {
	rc := context.Background()
	sk = db.Redis().Get(rc, ak).Val()

	if sk == "" {
		var aksk AKSK
		res := db.Mysql().Model(&AKSKExtension{}).Select("sk", "aid", "expire_time").Where("ak = ?", ak).First(&aksk)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				err = ErrNotFoundSecretKey
				return
			} else {
				err = res.Error
				etlog.L().Error(err.Error())
				return
			}
		}

		if aksk.ExpireTime > 0 {
			t := time.Since(time.Unix(aksk.ExpireTime, 0))
			if t > 0 {
				err = ErrAccessKeyExpired
				return
			}
		}

		sk = aksk.Sk
		db.Redis().Set(rc, ak, aksk.Sk, 0)
		if aksk.ExpireTime > 0 {
			db.Redis().ExpireAt(rc, ak, time.Unix(aksk.ExpireTime, 0))
		}
	}

	return
}

func ValidateSignature(ak string, sign string, content string) error {
	sk, err := getSK(ak)
	if err != nil {
		return err
	}

	signBytes, err := hex.DecodeString(sign)
	if err != nil {
		return err
	}

	hash := hmac.New(sha256.New, []byte(sk))
	hash.Write([]byte(content))
	ss := hash.Sum(nil)

	if ok := hmac.Equal(signBytes, ss); !ok {
		return ErrSignatureInvalid
	}

	return nil
}
