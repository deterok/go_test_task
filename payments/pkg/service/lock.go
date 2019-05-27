package service

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	redsync "gopkg.in/redsync.v1"
)

// LockRepository provides interface for simple distributed mechanism for restricting access to the same object
type LockRepository interface {
	Lock(key string) error
	Unlock(key string) error
}

type lockRepository struct {
	sync *redsync.Redsync
}

// NewLockReposytory resturns mechanism for restricting access using redis
func NewLockReposytory(redisAddr string) LockRepository {
	pool := &redis.Pool{
		MaxIdle:     1,
		IdleTimeout: 5 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	sync := redsync.New([]redsync.Pool{pool})

	return &lockRepository{sync}
}

func (l *lockRepository) Lock(key string) error {
	lock := l.sync.NewMutex(key)
	return errors.Wrap(lock.Lock(), "mutex locking failed")
}

func (l *lockRepository) Unlock(key string) error {
	lock := l.sync.NewMutex(key)
	if !lock.Unlock() {
		return errors.New("mutex releasing failed")
	}
	return nil
}
