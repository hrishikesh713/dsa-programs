package lru

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
)

func Setup() (*LRUCache, error) {
	size := 1000
	maxBytes := 1 << 20
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	c, err := NewLRUCache(size, WithMaxBytes(maxBytes), WithLogHandler(logHandler))
	if err != nil {
		return nil, fmt.Errorf("cannot initialize lrucache %with", err)
	}
	return c, nil
}

func TestLRUAdd(t *testing.T) {
	lruCache, err := Setup()
	if err != nil {
		t.Fatalf("cannot initialize lrucache %+v", err)
	}
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
	t.Run("find elements in the cache", func(t *testing.T) {
		v, err := lruCache.Get("foo")
		if err != nil {
			t.Errorf("should have found the foo entry %+v", err)
		}
		if v.Key != "foo" || v.Value != "bar" {
			t.Errorf("could not find the right entity : %+v", v)
		}
	})
	t.Run("evict the elements from the cache based on a TTL", func(t *testing.T) {
		_ = lruCache.Evict()
	})
}
