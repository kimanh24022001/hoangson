package database

import (
	"fmt"
	"testing"
)

func TestUpdateQuery(t *testing.T) {
	{
		qb := NewUpdateBuilder("users", []string{"id", "hello", "world"}, []any{"hello", "hello", "hello"})
		qb.WhereVar(true, "users", "id", Op_Eq, "HELLO")
		fmt.Printf("qb.String(): %v\n", qb.String())
	}

	{
		qb := NewSelectBuilder([]string{`"users"."id"`, `"users"."hello"`, `"users"."world"`})
		qb.From("users")
		qb.WhereVar(true, "users", "id", Op_Eq, "HELLO")
		qb.LimitVar(13)
		qb.OffsetVar(0)
		fmt.Printf("qb.String(): %v\n", qb.String())
	}
}
