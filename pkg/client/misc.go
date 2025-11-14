/*
Copyright 2024 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
)

func filterProps(in api.UntypedDto, props []string) api.UntypedDto {
	m := make(map[string]interface{})
	for k, v := range in {
		if !slices.Contains(props, k) {
			m[k] = v
		}
	}
	return m
}

func ensureResponseCode(r *http.Response, code int) error {
	if r.StatusCode != code {
		return fmt.Errorf("unexpected error code: %d, wanted: %d, body: %s",
			r.StatusCode, code, tryConsumeResponseBody(r))
	}
	return nil
}

func tryConsumeResponseBody(r *http.Response) string {
	var out bytes.Buffer
	if _, err := io.Copy(&out, r.Body); err != nil {
		return ""
	}
	return out.String()
}

type responseLoggingHttpDoer struct {
	d api.HttpRequestDoer
	l *slog.Logger
}

func (r *responseLoggingHttpDoer) Do(req *http.Request) (*http.Response, error) {
	resp, err := r.d.Do(req)
	if resp != nil {
		r.l.Debug("got response", "status", resp.StatusCode)
	}
	return resp, err
}

func ResponseLogger(d api.HttpRequestDoer, l *slog.Logger) api.HttpRequestDoer {
	return &responseLoggingHttpDoer{d: d, l: l}
}

func RequestLogger(l *slog.Logger) api.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		l.Debug("sending request", "method", req.Method, "url", req.URL.String())
		return nil
	}
}
