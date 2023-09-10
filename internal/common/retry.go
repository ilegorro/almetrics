package common

import (
	"context"
	"errors"
	"net/http"
	"syscall"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type funcExec func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
type funcDo func(req *http.Request) (*http.Response, error)

func WithRetryDo(aFunc funcDo, req *http.Request) (*http.Response, error) {
	var err error
	var resp *http.Response
	attempts := 0
	for {
		resp, err = aFunc(req)
		if attempts == 3 || !(errors.Is(err, syscall.ECONNREFUSED)) {
			break
		}
		attempts += 1
		time.Sleep(time.Duration((attempts*2)-1) * time.Second)
	}

	return resp, err
}

func WithRetryExec(aFunc funcExec, ctx context.Context, sql string, args ...any) error {
	var pgErr *pgconn.PgError
	var err error
	attempts := 0
	for {
		_, err = aFunc(ctx, sql, args...)
		if attempts == 3 || !(errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)) {
			break
		}
		attempts += 1
		time.Sleep(time.Duration((attempts*2)-1) * time.Second)
	}

	return err
}
