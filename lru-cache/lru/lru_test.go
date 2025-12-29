package lru

import (
	"container/list"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/hrishikesh713/dsa-programs/rate-limiter/utils"
)

func TestLRUAdd(t *testing.T) {
	size := 1000
	maxBytes := 1 << 20
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	lruCache, _ := NewLRUCache(size)
	t.Run("initialize lru with options and required args", func(t *testing.T) {
		c, err := NewLRUCache(size, WithMaxBytes(maxBytes), WithLogHandler(logHandler))
		if err != nil {
			t.Fatalf("cannot initialize Cache %#v due to err %#v", c, err)
		}
	})

	t.Run("add elements to lru cache", func(t *testing.T) {
		r := Entity{
			Key:       "foo",
			Value:     "bar",
			Timestamp: time.Now(),
		}
		err := lruCache.Add(r)
		if err != nil {
			t.Errorf("could not add value to the cache %w", err)
		}
	})

	t.Run("evict the elements from the cache based on a TTL", func(t *testing.T) {
		err := lruCache.Evict(time.Now().Add(-time.Hour))
	})
}


type Entity struct {
	Key       string
	Value     string
	Timestamp time.Time
}

func WithMaxBytes(maxBytes int) Opt {
	return func(lrucache *LRUCache) error {
		lrucache.Bytes = maxBytes
		return nil
	}
}

func WithLogHandler(logHandler slog.Handler) Opt {
	return func(lrucache *LRUCache) error {
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
	lc.DLL.PushBack(e)
	return nil
}
