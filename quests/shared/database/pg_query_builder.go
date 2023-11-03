package database

import (
	"strings"

	"smatyx.com/config"
)

type Statement int
type Operator int

const (
	Statement_Insert Statement = iota
	Statement_Update
	Statement_Select
	Statement_Delete
)

const (
	Op_Eq Operator = iota
	Op_Neq
	Op_Lt
	Op_Gt
	Op_Le
	Op_Ge
	Op_In
	Op_Like
	Op_Ilike
	Op_Is
)

var OperatorKeywords = []string{
	Op_Eq:    "=",
	Op_Neq:   "<>",
	Op_Lt:    "<",
	Op_Gt:    ">",
	Op_Le:    "<=",
	Op_Ge:    ">=",
	Op_In:    "IN",
	Op_Like:  "LIKE",
	Op_Ilike: "ILIKE",
	Op_Is:    "Is",
}

type QueryBuilder struct {
	Stmt Statement

	TableName string
	Columns   []string

	FromItems    []string
	WhereConds   []string
	OrderByExprs []string
	// TODO: having
	// TODO: join
	// TODO: group by
	Limit  string
	Offset string

	VarCount int
	Args     []any
}

func NewInsertBuilder(tableName string, columns []string, args []any) *QueryBuilder {
	if config.Debug {
		if len(args) % len(columns) != 0 {
			panic("")
		}
	}

	result := &QueryBuilder{
		Stmt:      Statement_Insert,
		TableName: tableName,
		Columns:   columns,
		Args:      args,
	}

	return result
}

func NewUpdateBuilder(tableName string, columns []string, args []any) *QueryBuilder {
	if config.Debug {
		if len(columns) != len(args) {
			panic("")
		}
	}

	result := &QueryBuilder{
		Stmt:       Statement_Update,
		TableName:  tableName,
		Columns:    columns,
		WhereConds: make([]string, 0, 10),
		VarCount:   0,
		Args:       make([]any, 0, len(args)+2),
	}
	result.Args = append(result.Args, args...)
	result.VarCount += len(args)

	return result
}

func NewSelectBuilder(columns []string) *QueryBuilder {
	result := &QueryBuilder{
		Stmt:         Statement_Select,
		TableName:    "",
		Columns:      columns,
		FromItems:    make([]string, 0, 3),
		WhereConds:   make([]string, 0, 10),
		OrderByExprs: make([]string, 0, 5),
		VarCount:     0,
	}

	return result
}

func (self *QueryBuilder) String() string {
	sb := strings.Builder{}

	switch self.Stmt {
	case Statement_Insert:
		sb.WriteString("INSERT INTO ")
		sbWritePgName(&sb, self.TableName)
		sb.WriteRune('(')
		sb.WriteString(strings.Join(self.Columns, ","))

		sb.WriteString(") VALUES")
		columnCount := len(self.Columns)
		rowCount := len(self.Args) / columnCount

		startIdx := 0
		for i := 0; i < rowCount; i++ {
			if i != 0 {
				sb.WriteRune(',')
			}
			endIdx := startIdx + columnCount
			variables := strings.Join(PresetVariables[startIdx:endIdx], ",")
			startIdx = endIdx

			sb.WriteRune('(')
			sb.WriteString(variables)
			sb.WriteRune(')')
		}
		sb.WriteRune(';')

	case Statement_Select:
		sb.WriteString("SELECT ")
		sb.WriteString(strings.Join(self.Columns, ","))

		sb.WriteString(" FROM ")
		for i, name := range self.FromItems {
			if i != 0 {
				sb.WriteRune(',')
			}
			sbWritePgName(&sb, name)
		}

		if len(self.WhereConds) > 0 {
			sb.WriteString(" WHERE")
			for _, cond := range self.WhereConds {
				sb.WriteRune(' ')
				sb.WriteString(cond)
			}
		}
		if len(self.OrderByExprs) > 0 {
			sb.WriteString(" ORDER BY")
			for i, cond := range self.OrderByExprs {
				if i != 0 {
					sb.WriteRune(',')
				}
				sb.WriteString(cond)
			}
		}
 
		if len(self.Limit) != 0 {
			sb.WriteString(" LIMIT ")
			sb.WriteString(self.Limit)
		}

		if len(self.Offset) != 0 {
			sb.WriteString(" OFFSET ")
			sb.WriteString(self.Offset)
		}

		sb.WriteRune(';')

	case Statement_Update:
		sb.WriteString("UPDATE ")
		sbWritePgName(&sb, self.TableName)
		sb.WriteString(" SET ")
		for i, column := range self.Columns {
			if i != 0 {
				sb.WriteRune(',')
			}
			sbWritePgName(&sb, column)
			sb.WriteRune('=')
			sb.WriteString(PresetVariables[i])
		}

		sb.WriteString(" WHERE")
		for _, cond := range self.WhereConds {
			sb.WriteRune(' ')
			sb.WriteString(cond)
		}
		sb.WriteRune(';')
	}

	return sb.String()
}

func (self *QueryBuilder) From(name string) {
	self.FromItems = append(self.FromItems, name)
}

func (self *QueryBuilder) WhereVar(table string, column string, op Operator, arg any) {
	sb := strings.Builder{}

	if len(self.WhereConds) != 0 {
		sb.WriteString("AND ")
	}

	sbWritePgName(&sb, table)
	sb.WriteRune('.')
	sbWritePgName(&sb, column)

	sb.WriteString(OperatorKeywords[op])
	sb.WriteString(PresetVariables[self.VarCount])

	self.WhereConds = append(self.WhereConds, sb.String())
	self.VarCount++

	self.Args = append(self.Args, arg)
}

func (self *QueryBuilder) WhereInVars(table, column string, args []any) {
	if len(args) == 0 {
		return
	}

	sb := strings.Builder{}

	if len(self.WhereConds) != 0{
		sb.WriteString("AND ")
	}

	varCount := len(args)
	values := strings.Join(PresetVariables[self.VarCount:self.VarCount+varCount], ",")

	sbWritePgName(&sb, table)
	sb.WriteRune('.')
	sbWritePgName(&sb, column)

	sb.WriteString(" IN (")
	sb.WriteString(values)
	sb.WriteRune(')')

	self.WhereConds = append(self.WhereConds, sb.String())
	self.VarCount += varCount

	self.Args = append(self.Args, args...)
}

func (self *QueryBuilder) OffsetVar(value int) {
	self.Offset = PresetVariables[self.VarCount]
	self.VarCount++

	self.Args = append(self.Args, value)
}

func (self *QueryBuilder) LimitVar(value int) {
	self.Limit = PresetVariables[self.VarCount]
	self.VarCount++

	self.Args = append(self.Args, value)
}

func (self *QueryBuilder) OrderBy(table, column string, isIncrease bool, nullFirst bool) {
	sb := strings.Builder{}
	sbWritePgName(&sb, table)
	sb.WriteRune('.')
	sbWritePgName(&sb, column)

	if isIncrease {
		sb.WriteString(" ASC")
	} else {
		sb.WriteString(" DESC")
	}

	if nullFirst {
		sb.WriteString(" NULLS FIRST")
	} else {
		sb.WriteString(" NULLS LAST")
	}

	self.OrderByExprs = append(self.OrderByExprs, sb.String())
}

func sbWritePgName(sb *strings.Builder, name string) {
	sb.WriteRune('"')
	sb.WriteString(strings.ReplaceAll(name, `"`, `""`))
	sb.WriteRune('"')
}
