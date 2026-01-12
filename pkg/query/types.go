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

package query

import "fmt"

type Op string

type Order interface {
	Name() string
	Asc() bool
}

type Orders []Order

type FilterExpression interface {
	fmt.Stringer
}

type Paging interface {
	fmt.Stringer
	Offset() uint64
	Size() int
}

type Interface interface {
	fmt.Stringer
	Orders() Orders
	Paging() Paging
	Filter() FilterExpression
}

type Builder interface {
	OrderBy(string, bool) Builder
	Paging(int, int) Builder
	Filter(FilterExpression) Builder
	Build() Interface
}

type InExpression interface {
	Name() string
	Values() []interface{}
}

type JunctionExpression interface {
	Op() Op
	Sub() []FilterExpression
}

type NotExpression interface {
	Sub() FilterExpression
}

type SimpleExpression interface {
	Name() string
	Op() Op
	Value() interface{}
}
