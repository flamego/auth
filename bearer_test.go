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
	"github.com/stretchr/testify/require"
)

func TestBearer(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		wantCode int
		wantBody string
	}{
		{
			name:     "good",
			token:    "foo",
			wantCode: http.StatusOK,
			wantBody: "foo",
		},
		{
			name:     "bad",
			token:    "bar",
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := flamego.NewWithLogger(&bytes.Buffer{})
			f.Use(Bearer("foo"))
			f.Get("/", func(token Token) string {
				return string(token)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", bearerPrefix+test.token)
			f.ServeHTTP(resp, req)

			assert.Equal(t, test.wantCode, resp.Code)
			assert.Equal(t, test.wantBody, resp.Body.String())
		})
	}
}

func TestBearerFunc(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		wantCode int
		wantBody string
	}{
		{
			name:     "primary token",
			header:   bearerPrefix + "foo",
			wantCode: http.StatusOK,
			wantBody: "foo",
		},
		{
			name:     "secondary token",
			header:   bearerPrefix + "bar",
			wantCode: http.StatusOK,
			wantBody: "bar",
		},
		{
			name:     "wrong token",
			header:   bearerPrefix + "nope",
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized\n",
		},
		{
			name:     "bad prefix",
			header:   "foo",
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := flamego.NewWithLogger(&bytes.Buffer{})
			f.Use(BearerFunc(func(token string) bool {
				return token == "foo" || token == "bar"
			}))
			f.Get("/", func(token Token) string {
				return string(token)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", test.header)
			f.ServeHTTP(resp, req)

			assert.Equal(t, test.wantCode, resp.Code)
			assert.Equal(t, test.wantBody, resp.Body.String())
		})
	}
}
