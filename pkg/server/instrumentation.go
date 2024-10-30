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

package server

import (
	"net/http"
)

type interceptedResp struct {
	delegate http.ResponseWriter
	written  int
	status   int
}

func (i *interceptedResp) Header() http.Header {
	return i.delegate.Header()
}

func (i *interceptedResp) Write(bytes []byte) (int, error) {
	size, err := i.delegate.Write(bytes)
	i.written += size
	return size, err
}

func (i *interceptedResp) WriteHeader(statusCode int) {
	i.delegate.WriteHeader(statusCode)
	i.status = statusCode
}

func (i *interceptedResp) Written() int {
	return i.written
}

func (i *interceptedResp) Status() int {
	if i.status == 0 {
		i.status = http.StatusOK
	}
	return i.status
}
