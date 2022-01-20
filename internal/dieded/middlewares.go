package dieded

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) CreateProfile(ctx context.Context, f ProfileForm) (p *Profile, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "CreateProfile", "name", f.Name, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.CreateProfile(ctx, f)
}

func (mw loggingMiddleware) GetProfile(ctx context.Context, id int) (p *Profile, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetProfile", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetProfile(ctx, id)
}

func (mw loggingMiddleware) QueryProfile(ctx context.Context, q Query) (p *Profile, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "QueryProfile", "name", q.Name, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.QueryProfile(ctx, q)
}

func (mw loggingMiddleware) DieProfile(ctx context.Context, id int) (p *Profile, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DieProfile", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DieProfile(ctx, id)
}

func (mw loggingMiddleware) DeleteProfile(ctx context.Context, id int) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteProfile", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteProfile(ctx, id)
}
