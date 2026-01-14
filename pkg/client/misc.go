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
	"net/http"
	"slices"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
	"github.com/rkosegi/db2rest-bridge/pkg/types"
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

// ensureResponseCode creates error if status code in HTTP response does not match expected one.
func ensureResponseCode(r *http.Response, code int, body []byte) error {
	if r.StatusCode != code {
		return &types.ErrorWithStatus{
			Status: code,
			Msg: fmt.Sprintf("unexpected error code: %d, wanted: %d, body: %s", r.StatusCode, code,
				string(bytes.TrimSpace(body)))}
	}
	return nil
}

func errorFromResponse(r *http.Response) error {
	return errorFromResponseWithMsg(r, fmt.Sprintf("unexpected error code: %d", r.StatusCode))
}

func errorFromResponseWithMsg(r *http.Response, msg string) error {
	return &types.ErrorWithStatus{
		Status: r.StatusCode,
		Msg:    msg,
	}
}
