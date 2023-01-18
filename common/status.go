//
// File: error_codes.go
// Created by Dizzrt on 2023/01/18.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package common

const (
	// auth
	StatusTokenUnhandled   = 1001
	StatusGenTokenError    = 1002
	StatusTokenMalformed   = 1003
	StatusTokenExpired     = 1004
	StatusTokenNotValidYet = 1005

	// user
	StatusUserUnhandled = 2001
	StatusEmailUsed     = 2002
	StatusRegisterOk    = 2003
	StatusUserNotExist  = 2004
	StatusWrongPassword = 2005
	StatusLoginOk       = 2006
)
