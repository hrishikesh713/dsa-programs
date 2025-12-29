package lru

import (
	"log/slog"
	"os"
	"testing"
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
			Key:   "foo",
			Value: "bar",
		}
		err := lruCache.Add(r)
		if err != nil {
			t.Errorf("could not add value to the cache %+v", err)
		}
	})

	t.Run("evict the elements from the cache based on a TTL", func(t *testing.T) {
		_ = lruCache.Evict()
	})
}
