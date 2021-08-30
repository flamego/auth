// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"crypto/sha512"
	"crypto/subtle"
)

// SecureCompare performs a constant time compare of two strings to prevent
// timing attacks.
func SecureCompare(given, actual string) bool {
	givenSHA := sha512.Sum512([]byte(given))
	actualSHA := sha512.Sum512([]byte(actual))
	return subtle.ConstantTimeCompare(givenSHA[:], actualSHA[:]) == 1
}
