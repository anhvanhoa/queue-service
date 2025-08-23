package database

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

// Hook to log the queries
type queryHook struct{}

func (h queryHook) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (h queryHook) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	bytes, _ := q.FormattedQuery()
	fmt.Println("After query\n" + string(bytes) + "\n")
	return nil
}

func NewQueryHook() pg.QueryHook {
	return queryHook{}
}
