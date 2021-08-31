// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flamego/flamego"
	"github.com/stretchr/testify/assert"
)

func TestBearer(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		wantCode int
	}{
		{
			name:     "good",
			token:    "foo",
			wantCode: http.StatusOK,
		},
		{
			name:     "bad",
			token:    "bar",
			wantCode: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := flamego.NewWithLogger(&bytes.Buffer{})
			f.Use(Bearer("foo"))
			f.Get("/", func() {})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			assert.Nil(t, err)

			req.Header.Set("Authorization", bearerPrefix+test.token)
			f.ServeHTTP(resp, req)

			assert.Equal(t, test.wantCode, resp.Code)
		})
	}
}

func TestBearerFunc(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		wantCode int
	}{
		{
			name:     "primary password",
			header:   bearerPrefix + "foo",
			wantCode: http.StatusOK,
		},
		{
			name:     "secondary password",
			header:   bearerPrefix + "bar",
			wantCode: http.StatusOK,
		},
		{
			name:     "wrong password",
			header:   bearerPrefix + "nope",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "bad prefix",
			header:   "foo",
			wantCode: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := flamego.NewWithLogger(&bytes.Buffer{})
			f.Use(BearerFunc(func(token string) bool {
				return token == "foo" || token == "bar"
			}))
			f.Get("/", func() {})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			assert.Nil(t, err)

			req.Header.Set("Authorization", test.header)
			f.ServeHTTP(resp, req)

			assert.Equal(t, test.wantCode, resp.Code)
		})
	}
}
