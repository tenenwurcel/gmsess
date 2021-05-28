package domain

import "context"

type RedisRepository interface {
	New(ctx context.Context, session Session) error
	Delete(ctx context.Context, sid string) error
	Fetch(ctx context.Context, sid string) (*Session, error)
}
