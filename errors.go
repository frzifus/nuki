// Copyright 2018 The Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package nuki

import (
	"errors"
)

var (
	// ErrNotImplemented function not implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrAuthFailed Auth() failed
	ErrAuthFailed = errors.New("authentication failed")
	// ErrInvalidURL : 400
	ErrInvalidURL = errors.New("given url is invalid or to long")
	// ErrInvalidToken : 401
	ErrInvalidToken = errors.New("token is invalid")
	// ErrAuthDisabled : 403
	ErrAuthDisabled = errors.New("authentication disabled")
	// ErrSmartLockUnknown : 404
	ErrSmartLockUnknown = errors.New("given smart lock is unknown")
	// ErrSmartLockOffline : 503
	ErrSmartLockOffline = errors.New("given smart lock is offline")
	// ErrUnknown : nobody knows
	ErrUnknown = errors.New("something went wrong")
)

// ErrorFromStatus returns the matching error based on the given status code
func ErrorFromStatus(status int) error {
	switch status {
	case 200:
		return nil
	case 400:
		return ErrInvalidURL
	case 401:
		return ErrInvalidToken
	case 403:
		return ErrAuthDisabled
	case 404:
		return ErrSmartLockUnknown
	case 503:
		return ErrSmartLockOffline
	default:
		return ErrUnknown
	}
}
