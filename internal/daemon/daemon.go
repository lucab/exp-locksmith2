package daemon

import (
	"context"
	"errors"
	"time"

	"go.etcd.io/etcd/clientv3"
	//"github.com/lucab/exp-locksmith2/internal/lock"
)

type lockManager struct {
	etcd3Client *clientv3.Client
}

func newLockManager(dialTimeout time.Duration) (*lockManager, error) {

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	manager := lockManager{
		etcd3Client: client,
	}

	return &manager, nil
}

func (lm *lockManager) recursiveLock(ctx context.Context, uuid string, group string) error {
	if lm == nil {
		return errors.New("nil lock manager")
	}

	return nil
}

func (lm *lockManager) unlockIfHeld(ctx context.Context, uuid string, group string) error {
	if lm == nil {
		return errors.New("nil lock manager")
	}

	return nil
}

func (lm *lockManager) Close() {
	if lm == nil {
		return
	}

	lm.etcd3Client.Close()
}
