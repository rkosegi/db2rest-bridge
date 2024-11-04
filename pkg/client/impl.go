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
	"encoding/json"
	"log/slog"

	"github.com/rkosegi/db2rest-bridge/pkg/api"
)

type generic[T any] struct {
	l     *slog.Logger
	be    string
	ent   string
	c     *api.ClientWithResponses
	decFn DecoderFn[T]
	encFn EncoderFn[T]
	// list of fields not to sent during updates
	skipProps []string
	// list of client options
	copts []api.ClientOption
}

func defaultEncoderFn[T any](obj *T) (api.UntypedDto, error) {
	var (
		err  error
		m    api.UntypedDto
		data []byte
	)
	if data, err = json.Marshal(obj); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func defaultDecoderFn[T any](m api.UntypedDto) (*T, error) {
	var (
		err  error
		t    T
		data []byte
	)
	if data, err = json.Marshal(m); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// SkipProperties sets properties that will not be sent to REST API during create/update operations
func SkipProperties[T any](props []string) Opt[T] {
	return func(g *generic[T]) {
		g.skipProps = props
	}
}

func WithPayloadEncoder[T any](encFn EncoderFn[T]) Opt[T] {
	return func(g *generic[T]) {
		g.encFn = encFn
	}
}

func WithPayloadDecoder[T any](decFn DecoderFn[T]) Opt[T] {
	return func(g *generic[T]) {
		g.decFn = decFn
	}
}

func WithLogger[T any](l *slog.Logger) Opt[T] {
	return func(g *generic[T]) {
		g.l = l
	}
}

func WithClientOptions[T any](copts ...api.ClientOption) Opt[T] {
	return func(g *generic[T]) {
		g.copts = copts
	}
}

func New[T any](endpoint, backend, entity string, opts ...Opt[T]) (GenericInterface[T], error) {
	g := &generic[T]{
		be:  backend,
		ent: entity,
	}
	for _, opt := range append([]Opt[T]{
		WithLogger[T](slog.Default()),
		WithPayloadDecoder[T](defaultDecoderFn[T]),
		WithPayloadEncoder[T](defaultEncoderFn[T]),
		SkipProperties[T]([]string{}),
	}, opts...) {
		opt(g)
	}
	c, err := api.NewClientWithResponses(endpoint, g.copts...)
	if err != nil {
		return nil, err
	}
	g.c = c
	g.l = g.l.With("backend", g.be, "entity", g.ent)
	return g, nil
}
