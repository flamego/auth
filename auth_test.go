// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecureCompare(t *testing.T) {
	tests := []struct {
		given  string
		actual string
		want   bool
	}{
		{"foo", "foo", true},
		{"bar", "bar", true},
		{"password", "password", true},
		{"Foo", "foo", false},
		{"foo", "foobar", false},
		{"password", "pass", false},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			got := SecureCompare(test.given, test.actual)
			assert.Equal(t, test.want, got)
		})
	}
}
