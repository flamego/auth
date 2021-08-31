// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/flamego/flamego"
)

const basicPrefix = "Basic "

// Basic returns a middleware handler that injects auth.User into the request
// context upon successful basic authentication. The handler responds
// http.StatusUnauthorized when authentication fails.
func Basic(username, password string) flamego.Handler {
	want := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return flamego.ContextInvoker(func(c flamego.Context) {
		got := c.Request().Header.Get("Authorization")
		if !SecureCompare(basicPrefix+want, got) {
			basicUnauthorized(c.ResponseWriter())
			return
		}
		c.Map(User(username))
	})
}

// BasicFunc returns a middleware handler that injects auth.User into the
// request context upon successful basic authentication with the given function.
// The function should return true for a valid username and password
// combination.
func BasicFunc(fn func(username, password string) bool) flamego.Handler {
	return flamego.ContextInvoker(func(c flamego.Context) {
		auth := c.Request().Header.Get("Authorization")
		n := len(basicPrefix)
		if len(auth) < n || auth[:n] != basicPrefix {
			basicUnauthorized(c.ResponseWriter())
			return
		}
		b, err := base64.StdEncoding.DecodeString(auth[n:])
		if err != nil {
			basicUnauthorized(c.ResponseWriter())
			return
		}
		tokens := strings.SplitN(string(b), ":", 2)
		if len(tokens) != 2 || !fn(tokens[0], tokens[1]) {
			basicUnauthorized(c.ResponseWriter())
			return
		}
		c.Map(User(tokens[0]))
	})
}

func basicUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
