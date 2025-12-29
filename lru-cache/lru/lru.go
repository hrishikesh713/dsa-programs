package lru

import (
	"container/list"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/hrishikesh713/dsa-programs/rate-limiter/utils"
)

type Entity struct {
	Key   string
	Value string
}

func WithMaxBytes(maxBytes int) Opt {
	return func(lrucache *LRUCache) error {
		lrucache.Bytes = maxBytes
		return nil
	}
}

func WithLogHandler(logHandler slog.Handler) Opt {
	return func(lrucache *LRUCache) error {
		switch v := logHandler.(type) {
		case *slog.JSONHandler:
			if v == nil {
				return fmt.Errorf("pointer inside loghandler is nil")
			}
		default:
			if v == nil {
				return fmt.Errorf("interface is nil")
			}
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Pointer && rv.IsNil() {
				return fmt.Errorf("cannot have a nil pointer stored in interface")
			}
			return fmt.Errorf("value inside loghandler interface is not one that is recognized %T", v)
		}
		lrucache.Logger = slog.New(logHandler)
		return nil
	}
}

type Opt func(lrucache *LRUCache) error

func WithClock(clock utils.Clock) Opt {
	return func(lrucache *LRUCache) error {
		lrucache.Clock = clock
		return nil
	}
}

type LRUCache struct {
	rwmu   sync.RWMutex
	DLL    *list.List
	Clock  utils.Clock
	Logger *slog.Logger
	Size   int
	Bytes  int
}

func NewLRUCache(size int, opts ...Opt) (*LRUCache, error) {
	lrucache := &LRUCache{
		DLL:    list.New(),
		Clock:  utils.NewFakeClock(),
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
	lc.rwmu.Lock()
	defer lc.rwmu.Unlock()
	lc.Logger.Info("pushing entries into cache")
	lc.DLL.PushBack(e)
	return nil
}

func (lc *LRUCache) Evict() error {
	lc.rwmu.Lock()
	defer lc.rwmu.Unlock()
	lc.Logger.Info("evicting entries")
	if lc.DLL.Len() > 0 {
		lc.DLL.Remove(lc.DLL.Front())
	}
	return nil
}

func (lc *LRUCache) Get(key string) (Entity, error) {
	lc.rwmu.RLock()
	defer lc.rwmu.RUnlock()
	for e := lc.DLL.Back(); e != nil; e = e.Prev() {
		if v, ok := e.Value.(Entity); ok && v.Key == key {
			return v, nil
		}
	}
	return Entity{}, fmt.Errorf("cannot find the value")
}
