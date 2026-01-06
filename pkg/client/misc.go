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
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/samber/lo"
)

// excludeProps return untyped object with all provided properties removed
func excludeProps(in api.UntypedDto, props []string) api.UntypedDto {
	return lo.OmitBy(in, func(key string, value interface{}) bool {
		return slices.Contains(props, key)
	})
}

// onlyProps return untyped object with only provided properties
func onlyProps(in api.UntypedDto, props []string) api.UntypedDto {
	return lo.OmitBy(in, func(key string, value interface{}) bool {
		return !slices.Contains(props, key)
	})
}

func ensureResponseCode(r *http.Response, code int) error {
	if r.StatusCode != code {
		return fmt.Errorf("unexpected error code: %d, wanted: %d, body: %s",
			r.StatusCode, code, tryConsumeResponseBody(r))
	}
	return nil
}

func invalidCode(r *http.Response) error {
	return fmt.Errorf("unexpected error code: %d: %s", r.StatusCode, tryConsumeResponseBody(r))
}

func tryConsumeResponseBody(r *http.Response) string {
	var out bytes.Buffer
	if _, err := io.Copy(&out, r.Body); err != nil {
		return ""
	}
	return out.String()
}
