package kv

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	ckv "github.com/hardcore-os/corekv"
	"github.com/hardcore-os/corekv/utils"
)

type KV interface {
	Get(ctx context.Context, key string) (value []byte, ok bool, err error)
	Set(ctx context.Context, key string, value []byte, ttlSec int) error
	Delete(ctx context.Context, key string) error
}

type CoreKV struct {
	path string
	db   ckv.CoreAPI
	mu   sync.RWMutex
}

func NewCoreKV(path string) (*CoreKV, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("corekv data path is required")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve corekv path: %w", err)
	}
	if err := os.MkdirAll(absPath, 0o755); err != nil {
		return nil, fmt.Errorf("create corekv dir: %w", err)
	}

	opts := ckv.NewDefaultOptions()
	opts.WorkDir = absPath
	opts.MemTableSize = 32 << 20
	opts.SSTableMaxSz = 256 << 20
	opts.ValueLogFileSize = 64 << 20

	db := ckv.Open(opts)
	if db == nil {
		return nil, errors.New("corekv open returned nil")
	}

	return &CoreKV{path: absPath, db: db}, nil
}

func (c *CoreKV) Get(ctx context.Context, key string) ([]byte, bool, error) {
	if c == nil || c.db == nil {
		return nil, false, errors.New("corekv is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	if len(key) == 0 {
		return nil, false, errors.New("key is required")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, err := c.db.Get([]byte(key))
	if err != nil {
		if errors.Is(err, utils.ErrKeyNotFound) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("corekv get %s: %w", key, err)
	}
	if entry == nil || entry.IsDeletedOrExpired() {
		return nil, false, nil
	}

	value := make([]byte, len(entry.Value))
	copy(value, entry.Value)
	return value, true, nil
}

func (c *CoreKV) Set(ctx context.Context, key string, value []byte, ttlSec int) error {
	if c == nil || c.db == nil {
		return errors.New("corekv is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(key) == 0 {
		return errors.New("key is required")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	entry := utils.NewEntry([]byte(key), append([]byte(nil), value...))
	if ttlSec > 0 {
		entry.WithTTL(time.Duration(ttlSec) * time.Second)
	}

	if err := c.db.Set(entry); err != nil {
		return fmt.Errorf("corekv set %s: %w", key, err)
	}
	return nil
}

func (c *CoreKV) Delete(ctx context.Context, key string) error {
	if c == nil || c.db == nil {
		return errors.New("corekv is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(key) == 0 {
		return errors.New("key is required")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.db.Del([]byte(key)); err != nil {
		if errors.Is(err, utils.ErrKeyNotFound) {
			return nil
		}
		return fmt.Errorf("corekv delete %s: %w", key, err)
	}
	return nil
}

func (c *CoreKV) Close() error {
	if c == nil || c.db == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.db.Close()
	c.db = nil
	return err
}
