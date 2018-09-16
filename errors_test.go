// Copyright 2018 The Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package nuki

import "testing"

func TestErrorFromStatus(t *testing.T) {
	tt := []struct {
		name     string
		status   int
		expected error
	}{
		{
			name:     "status OK",
			status:   200,
			expected: nil,
		},
		{
			name:     "given url is invalid or to long",
			status:   400,
			expected: ErrInvalidURL,
		},
		{
			name:     "token is invalid",
			status:   401,
			expected: ErrInvalidToken,
		},
		{
			name:     "authentication disabled",
			status:   403,
			expected: ErrAuthDisabled,
		},
		{
			name:     "given smart lock is unknown",
			status:   404,
			expected: ErrSmartLockUnknown,
		},
		{
			name:     "given smart lock is offline",
			status:   503,
			expected: ErrSmartLockOffline,
		},
		{
			name:     "unknown code",
			status:   42,
			expected: ErrUnknown,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if err := ErrorFromStatus(tc.status); err != tc.expected {
				t.Errorf("ErrorFromStatus() error = %v, expected %v", err, tc.expected)
			}
		})
	}
}
