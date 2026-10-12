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

package types

import (
	"net/http"
)

func NewErrorWithStatus(msg string, status int) error {
	return &ErrorWithStatus{
		Status: status,
		Msg:    msg,
	}
}

func WrapError(msg string, err error) error {
	return &ErrorWithStatus{
		Status:  http.StatusInternalServerError,
		Msg:     msg,
		Wrapped: err,
	}
}

func WrapErrorWithStatus(msg string, err error, status int) error {
	return &ErrorWithStatus{
		Status:  status,
		Msg:     msg,
		Wrapped: err,
	}
}

type ErrorWithStatus struct {
	Status  int
	Msg     string
	Wrapped error
}

func (h *ErrorWithStatus) Unwrap() error {
	return h.Wrapped
}

func (h *ErrorWithStatus) Error() string {
	return h.Msg
}
