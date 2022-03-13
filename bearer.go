// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"net/http"

	"github.com/flamego/flamego"
)

var bearerPrefix = "Bearer "

// Token is the authenticated token that was extracted from the request.
type Token string

// Bearer returns a middleware handler that injects auth.User (empty string)
// into the request context upon successful bearer authentication. The handler
// responds http.StatusUnauthorized when authentication fails.
func Bearer(token string) flamego.Handler {
	return flamego.ContextInvoker(func(c flamego.Context) {
		got := c.Request().Header.Get("Authorization")
		if !SecureCompare(bearerPrefix+token, got) {
			bearerUnauthorized(c.ResponseWriter())
			return
		}
		c.Map(Token(token))
	})
}

// BearerFunc returns a middleware handler that injects auth.User (empty string)
// into the request context upon successful bearer authentication with the given
// function. The function should return true for a valid bearer token.
func BearerFunc(fn func(token string) bool) flamego.Handler {
	return flamego.ContextInvoker(func(c flamego.Context) {
		auth := c.Request().Header.Get("Authorization")
		n := len(bearerPrefix)
		if len(auth) < n || auth[:n] != bearerPrefix {
			bearerUnauthorized(c.ResponseWriter())
			return
		}

		token := auth[n:]
		if !fn(token) {
			bearerUnauthorized(c.ResponseWriter())
			return
		}
		c.Map(Token(token))
	})
}

func bearerUnauthorized(res http.ResponseWriter) {
	http.Error(res, "Unauthorized", http.StatusUnauthorized)
}
