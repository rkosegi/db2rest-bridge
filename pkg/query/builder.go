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

type builder struct {
	ords Orders
	fe   FilterExpression
	pg   *page
}

func (b *builder) OrderBy(name string, asc bool) Builder {
	b.ords = append(b.ords, OrderBy(name, asc))
	return b
}

func (b *builder) Paging(offset uint64, size int) Builder {
	b.pg.offset = offset
	b.pg.size = size
	return b
}

func (b *builder) Filter(fe FilterExpression) Builder {
	b.fe = fe
	return b
}

func (b *builder) Build() Interface {
	return &qryData{orders: b.ords, paging: b.pg, filter: b.fe}
}

func BuilderFrom(qry Interface) Builder {
	return &builder{
		ords: qry.(*qryData).orders,
		pg:   qry.(*qryData).paging.(*page),
		fe:   qry.(*qryData).filter,
	}
}

func NewBuilder() Builder {
	return &builder{
		ords: Orders{},
		pg:   &page{offset: DefaultPageOffset, size: DefaultPageSize},
	}
}
