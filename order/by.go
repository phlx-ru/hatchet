package order

import (
	"entgo.io/ent/dialect/sql"
)

type By struct {
	Desc  bool
	Field string
}

func (b *By) OrderFunc() func(*sql.Selector) {
	if b.Desc {
		return func(s *sql.Selector) {
			s.OrderBy(sql.Desc(s.C(b.Field)))
		}
	}
	return func(s *sql.Selector) {
		s.OrderBy(sql.Asc(s.C(b.Field)))
	}
}

func ByAsc(field string) *By {
	return &By{
		Desc:  false,
		Field: field,
	}
}

func ByDesc(field string) *By {
	return &By{
		Desc:  true,
		Field: field,
	}
}

func ByFromRequest(field string, direction int32) *By {
	if direction == 0 {
		return ByAsc(field)
	}
	return ByDesc(field)
}
