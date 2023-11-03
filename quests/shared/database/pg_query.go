package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"smatyx.com/config"
)

const (
	QueryModeRead = iota
	QueryModeWrite
)

type Query struct {
	Mode       int
	Sql        string
	Args       []any
	ProcessRep ProcessRepFunc

	Done    chan struct{}
	Error   error
	Context context.Context
}

type Queries struct {
	Len  int
	Mode []int

	Sql  []string
	Args [][]any

	ProcessReps []ProcessRepFunc
	Errors      []error

	Done    chan bool
	Context context.Context
}

func NewReadQuery(
	context context.Context,
	scanRow ProcessRepFunc,
	sql string,
	arguments ...any) *Query {
	result := &Query{
		Mode:       QueryModeRead,
		Sql:        sql,
		Args:       arguments,
		ProcessRep: scanRow,
		Done:       make(chan struct{}),
		Context:    context,
	}

	return result
}

func NewWriteQuery(
	context context.Context,
	validateRep ProcessRepFunc,
	sql string,
	arguments ...any) *Query {
	result := &Query{
		Mode:       QueryModeWrite,
		Sql:        sql,
		Args:       arguments,
		ProcessRep: validateRep,
		Done:       make(chan struct{}),
		Context:    context,
	}

	return result
}

func (query *Query) Submit() error {
	t0 := time.Now()
	defer logElapsedTime(t0)

	if config.Debug {
		sb := strings.Builder{}
		sb.WriteRune('\n')
		for i, arg := range query.Args {
			sb.WriteRune('\t')
			sb.WriteString(fmt.Sprintf("%v %*v %*T | %v",
				PresetVariables[i],
				4-len(PresetVariables[i]), "|",
				13, arg, arg))
			sb.WriteRune('\n')
		}

		log.Printf("Submit new query:\nText: %v\nArgs: %v\n",
			query.Sql,
			sb.String())
	}

	if query.Mode == QueryModeRead {
		// return DbRead(query.Context, query.ProcessRep, query.Sql, query.Args...)
		return DbReadQuery(query)
	} else {
		return DbWrite(query.Context, query.ProcessRep, query.Sql, query.Args...)
	}
}

func NewQueries(ctx context.Context, capacity int) *Queries {
	result := &Queries{
		Len:          0,
		Mode:         make([]int, 0, capacity),
		Sql:          make([]string, 0, capacity),
		Args:         make([][]any, 0, capacity),
		ProcessReps:  make([]ProcessRepFunc, 0, capacity),
		Errors:       make([]error, 0, capacity),
		Done:         make(chan bool, 1),
		Context:      ctx,
	}
	return result
}

func NewQueriesFromBuilders(
	ctx context.Context,
	builders []*QueryBuilder,
	processReps []ProcessRepFunc) *Queries {

	if config.Debug {
		if len(builders) != len(processReps) {
			panic("The len of builders should be equal to processReps.")
		}
	}

	queriesLen := len(builders)
	result := &Queries{
		Len:          queriesLen,
		Mode:         make([]int, queriesLen),
		Sql:          make([]string, queriesLen),
		Args:         make([][]any, queriesLen),
		ProcessReps:  processReps,
		Errors:       make([]error, queriesLen),
		Done:         make(chan bool, 1),
		Context:      ctx,
	}

	for i := 0; i < queriesLen; i++ {
		qb := builders[i]
		if qb.Stmt == Statement_Select {
			result.Mode[i] = QueryModeRead
		} else {
			result.Mode[i] = QueryModeWrite
		}

		result.Sql[i] = qb.String()
		result.Args[i] = qb.Args
	}

	return result	
}

func (queries *Queries) AddReadQuery(
	scanRow ProcessRepFunc,
	sql string,
	arguments ...any) {
	// NOTE: We do not need to check the boundary here.

	queries.Mode = append(queries.Mode, QueryModeRead)
	queries.Sql = append(queries.Sql, sql)
	queries.Args = append(queries.Args, arguments)
	queries.ProcessReps = append(queries.ProcessReps, scanRow)
	queries.Errors = append(queries.Errors, nil)

	queries.Len++
}

func (queries *Queries) AddWriteQuery(
	validateRep ProcessRepFunc,
	sql string,
	arguments ...any) {
	// NOTE: We do not need to check the boundary here.

	queries.Mode = append(queries.Mode, QueryModeRead)
	queries.Sql = append(queries.Sql, sql)
	queries.Args = append(queries.Args, arguments)
	queries.ProcessReps = append(queries.ProcessReps, validateRep)
	queries.Errors = append(queries.Errors, nil)

	queries.Len++
}

func (queries *Queries) Submit() []error {
	t0 := time.Now()
	defer logElapsedTime(t0)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(queries.Len)

	for i := 0; i < queries.Len; i++ {
		mode := queries.Mode[i]
		text := queries.Sql[i]
		args := queries.Args[i]

		if config.Debug {
			sb := strings.Builder{}
			sb.WriteRune('\n')
			for i, arg := range args {
				sb.WriteRune('\t')
				sb.WriteString(fmt.Sprintf("%v %*v %*T | %v",
					PresetVariables[i],
					4-len(PresetVariables[i]), "|",
					13, arg, arg))
				sb.WriteRune('\n')
			}

			log.Printf("Submit new query:\nText: %v\nArgs: %v\n",
				text,
				sb.String())
		}

		if mode == QueryModeRead {
			idx := i
			go func() {
				// queries.Errors[idx] = DbRead(
				// 		queries.Context,
				// 		queries.ProcessReps[idx],
				// 		text, args...)

				// STUDY(duong): This strategy may be terrible for
				// local testing, but it works well when we have
				// networking i/o.
				queries.Errors[idx] = DbReadQuery(
					NewReadQuery(
						queries.Context,
						queries.ProcessReps[idx],
						text, args...))

				waitGroup.Done()
			}()
		} else {
			idx := i
			go func() {
				defer waitGroup.Done()

				queries.Errors[idx] = DbWrite(queries.Context, queries.ProcessReps[idx], text, args...)
			}()
		}
	}

	waitGroup.Wait()

	return queries.Errors
}

func logElapsedTime(t0 time.Time) {
	if config.Debug {
		log.Printf("Taken %d ms", time.Now().Sub(t0).Milliseconds())
	}
}
