// NOTE(duong): Due to my laziness and limited knowledge, we are
// currently using a clunky postgres driver.
//
// Our hope for the future is to be able to have a better driver or,
// even better, write our own database.
package database

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRep struct {
	Scanner    pgx.Rows
	CommandTag pgconn.CommandTag
}

func (self *PgRep) Scan(dest ...any) error {
	return self.Scanner.Scan(dest...)
}

func (self *PgRep) RowsAffected() int {
	return int(self.CommandTag.RowsAffected())
}

var Pool *pgxpool.Pool

func InitPg(dbUrl string) {
	cfg, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		panic(err)
	}
	fmt.Printf("cfg.MaxConns: %v\n", cfg.MaxConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	Pool = pool

	go Consumer()
}

type ProcessRepFunc func(result *PgRep) error

func DbReadQuery(query *Query) error {
	QueryChannel <- query

	select {
	case <-query.Done:
		return query.Error
	case <-time.After(1 * time.Second):
		return errors.New("Db read timeout")
	}
}

func DbRead(
	ctx context.Context,
	scanRow ProcessRepFunc,
	text string,
	args ...any) error {
	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, text, args...)
	if err != nil {
		return err
	}

	rep := PgRep{Scanner: rows}
	for rows.Next() {
		if scanRow != nil {
			err := scanRow(&rep)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DbWrite(
	ctx context.Context,
	validateRep ProcessRepFunc,
	text string,
	args ...any) error {
	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	commandTag, err := conn.Exec(ctx, text, args...)
	if err != nil {
		return err
	}
	rep := PgRep{CommandTag: commandTag}

	if validateRep != nil {
		return validateRep(&rep)
	}
	return nil
}

var QueryChannel chan *Query

func Consumer() int {
	var QueryChannelCap = int(Pool.Config().MaxConns * 100)
	QueryChannel = make(chan *Query, QueryChannelCap)

	queries := make([]*Query, 0, QueryChannelCap)
	for {
		queries = queries[:0] // NOTE(duong): clear the slice

		ctx := context.Background()

		queries = append(queries, <-QueryChannel)
		queriesLen := min(len(QueryChannel), QueryChannelCap - 1) + 1

		for i := 1; i < queriesLen; i++ {
			queries = append(queries, <-QueryChannel)
		}

		// NOTE(Duong): Start combine and execute read queries.
		//              First, put queries in buckets. Then request
		//              in multiple threads and connections.
		bucketCount := int(Pool.Config().MaxConns)
		if queriesLen < bucketCount {
			bucketCount = queriesLen
		}
		itemsPerBucket := queriesLen / bucketCount
		leftoverCount := queriesLen % bucketCount

		var waitGroup sync.WaitGroup
		waitGroup.Add(bucketCount)

		start := 0
		bound := itemsPerBucket

		for i := 0; i < bucketCount; i++ {
			if i < leftoverCount {
				bound++
			}

			subQueries := queries[start:bound]
			go func() {
				conn, err := Pool.Acquire(ctx)

				// NOTE(duong): for now, just panic
				if err != nil {
					for _, q := range subQueries {
						q.Error = err
						close(q.Done)
					}

					return
				}
				defer conn.Release()

				batch := pgx.Batch{}
				for _, q := range subQueries {
					batch.Queue(q.Sql, q.Args...)
				}

				results := conn.Conn().SendBatch(ctx, &batch)
				rep := PgRep{}

				for _, q := range subQueries {
					rep.Scanner, err = results.Query()

					if err != nil {
						q.Error = err
						close(q.Done)
						continue
					}
					if q.ProcessRep != nil {
						for rep.Scanner.Next() {
							err = q.ProcessRep(&rep)
							if err != nil {
								q.Error = err
								break
							}
						}
					}
					close(q.Done)
				}
				
				waitGroup.Done()
			}()

			start = bound
			bound = start + itemsPerBucket
		}
		waitGroup.Wait()
	}
}
