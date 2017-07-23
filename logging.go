package usermanagementsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
)

type middleware func(Service) Service

func ServiceLoggingMiddleware(logger log.Logger) middleware {
	return func(next Service) Service {
		return loggingMiddlware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddlware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddlware) CreateUser(ctx context.Context, email, password string) (err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "CreateUser",
			"email", email,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.CreateUser(ctx, email, password)
	return
}
