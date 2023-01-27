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
	// common
	StatusCommonOK    = 1001
	StatusCommonError = 1002

	// auth
	StatusTokenUnhandled   = 2001
	StatusGenTokenError    = 2002
	StatusTokenMalformed   = 2003
	StatusTokenExpired     = 2004
	StatusTokenNotValidYet = 2005

	// user
	StatusUserUnhandled     = 3001
	StatusEmailUsed         = 3002
	StatusRegisterOk        = 3003
	StatusUserNotExist      = 3004
	StatusWrongPassword     = 3005
	StatusLoginOk           = 3006
	StatusUserOk            = 3007
	StatusUserError         = 3008
	StatusUserLowPermission = 3009
)
