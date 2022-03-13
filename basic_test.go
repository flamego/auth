// Copyright 2021 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flamego/flamego"
)

func TestBasic(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantCode int
		wantBody string
	}{
		{
			name:     "good",
			username: "foo",
			password: "bar",
			wantCode: http.StatusOK,
			wantBody: "foo",
		},
		{
			name:     "bad",
			username: "bar",
			password: "foo",
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := flamego.NewWithLogger(&bytes.Buffer{})
			f.Use(Basic("foo", "bar"))
			f.Get("/", func(user User) string {
				return string(user)
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			auth := strings.Join([]string{test.username, test.password}, ":")
			req.Header.Set("Authorization", basicPrefix+base64.StdEncoding.EncodeToString([]byte(auth)))
			f.ServeHTTP(resp, req)

			assert.Equal(t, test.wantCode, resp.Code)
			assert.Equal(t, test.wantBody, resp.Body.String())
		})
	}
}

func TestBasicFunc(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		wantCode int
		wantBody string
	}{
		{
			name:     "primary password",
			header:   basicPrefix + "Zm9vOmJhcg==", // foo:bar
			wantCode: http.StatusOK,
			wantBody: "foo",
		},
		{
			name:     "secondary password",
			header:   basicPrefix + "Zm9vOmJheg==", // foo:baz
			wantCode: http.StatusOK,
			wantBody: "foo",
		},
		{
			name:     "wrong password",
			header:   basicPrefix + "Zm9vOm5vcGU=", // foo:nope
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized\n",
		},
		{
			name:     "bad prefix",
			header:   "Zm9vOmJheg==", // foo:baz
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized\n",
		},
		{
			name:     "bad encoding",
			header:   basicPrefix + "Zm9vOm5",
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := flamego.NewWithLogger(&bytes.Buffer{})
			f.Use(BasicFunc(func(username, password string) bool {
				return username == "foo" && (password == "bar" || password == "baz")
			}))
			f.Get("/", func(user User) string {
				return string(user)
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
