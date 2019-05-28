package service

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	redsync "gopkg.in/redsync.v1"
)

// Lock provides interface for simple distributed mechanism for restricting access to the same object
type Lock interface {
	Lock() error
	Unlock() error
}

type LockFactory interface {
	Make(key string) Lock
}

// ─── LOCK IMPLEMENTATION ────────────────────────────────────────────────────────

type lock struct {
	m *redsync.Mutex
}

func NewLock(m *redsync.Mutex) Lock {
	return &lock{m}
}

func (l *lock) Lock() error {
	if err := l.m.Lock(); err != nil {
		return errors.Wrap(err, "mutex locking failed")
	}

	return nil
}

func (l *lock) Unlock() error {
	if !l.m.Unlock() {
		return errors.New("mutex releasing failed")
	}

	return nil
}

// ─── LOCK POOL IMPLEMENTATION ───────────────────────────────────────────────────

type lockPool struct {
	locks []Lock
}

func NewLockPool(locks []Lock) Lock {
	return &lockPool{locks}
}

func (p *lockPool) Lock() error {
	for _, l := range p.locks {
		if err := l.Lock(); err != nil {
			return err
		}
	}

	return nil
}

func (p *lockPool) Unlock() error {
	for _, l := range p.locks {
		if err := l.Unlock(); err != nil {
			return err
		}
	}

	return nil
}

// ─── LOCK FACTORY IMPLEMENTATION ────────────────────────────────────────────────

type lockFactory struct {
	s *redsync.Redsync
}

func NewLockFactory(pool *redis.Pool) LockFactory {
	return &lockFactory{
		s: redsync.New([]redsync.Pool{pool}),
	}
}

func (f *lockFactory) Make(key string) Lock {
	return NewLock(f.s.NewMutex(key))
}
