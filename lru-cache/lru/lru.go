// Package lru is a LRUCache that has following properties:
// 1. Read heavy operations
// 2. Evict only when cache is full
// 3. Evict only the oldest entries
// 4. Oldest entries are decided based on how long they are in the cache
// 5. There are no duplicates.
// 6. If duplicate entry is inserted its value is updated.
// 7. Cache is an in-memory key value store.
// 8. key is a string and value is of type Any
// Limitations of LRUCache:
// 1. Currently the cache is single threaded
// 2. You cannot expand the size of the cache dynamically
// API of LRUCache:
// 1. Add to cache
// 2. Evict from cache
// 3. Initialize the cache
// 4. Clear the cache
// Edge cases ( what happens when?):
// 1. cache is initalized with zero capacity
// 2. cache is cleared and then evict is called.
// 3. empty key is inserted
// 4. same key inserted multiple times. (DDoS)
package lru

import (
	"container/list"
	"fmt"
	"log/slog"

	"github.com/hrishikesh713/dsa-programs/utils"
)

type Entity struct {
	Key   string
	Value any
}

func WithLogHandler(logHandler slog.Handler) Opt {
	return func(lrucache *LRUCache) error {
		switch v := logHandler.(type) {
		case *slog.JSONHandler:
			if v == nil {
				return fmt.Errorf("pointer inside loghandler is nil")
			}
		case *slog.TextHandler:
			if v == nil {
				return fmt.Errorf("pointer inside loghandler is nil")
			}
		default:
			r, err := utils.IsValueNil(v)
			if err != nil {
				return fmt.Errorf("wrong interface value %w", err)
			}
			if r {
				return fmt.Errorf("pointer inside interface is nil")
			}
			return fmt.Errorf("value inside loghandler interface is not one that is recognized %T", v)
		}
		lrucache.Logger = slog.New(logHandler)
		return nil
	}
}

type Opt func(lrucache *LRUCache) error

type LRUCache struct {
	DLL    *list.List
	Lookup map[string]*list.Element
	Logger *slog.Logger
	Size   int
}

func NewLRUCache(size int, opts ...Opt) (*LRUCache, error) {
	lrucache := &LRUCache{
		DLL:    list.New(),
		Lookup: make(map[string]*list.Element),
		Logger: slog.Default(),
		Size:   size,
	}
	for _, opt := range opts {
		if err := opt(lrucache); err != nil {
			return nil, fmt.Errorf("error while initializing %w", err)
		}
	}

	return lrucache, nil
}

func (lc *LRUCache) Add(e Entity) error {
	lc.Logger.Info("pushing entries into cache")
	if lc.DLL.Len() >= lc.Size {
		return fmt.Errorf("no capacity")
	}
	if v, ok := lc.Lookup[e.Key]; ok {
		v.Value = &e
		lc.DLL.MoveToBack(v)
		return nil
	}
	ee := lc.DLL.PushBack(&e)
	lc.Lookup[e.Key] = ee
	return nil
}

func (lc *LRUCache) Evict() error {
	lc.Logger.Info("evicting entries")
	if lc.DLL.Len() > 0 {
		e := lc.DLL.Front()
		delete(lc.Lookup, e.Value.(*Entity).Key)
		lc.DLL.Remove(e)
	}
	return nil
}

func (lc *LRUCache) Get(key string) (Entity, error) {
	if e, ok := lc.Lookup[key]; ok {
		lc.DLL.MoveToBack(e)
		return *e.Value.(*Entity), nil
	}
	return Entity{}, fmt.Errorf("cannot find the value")
}
