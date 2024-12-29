package jankdb_test

import (
	"testing"
	"time"

	"github.com/guarzo/jankdb"
)

func TestCache_SetGet(t *testing.T) {
	cache := jankdb.NewCache[string](5*time.Minute, 1*time.Minute)

	cache.Set("mykey", "myval")
	val, found := cache.Get("mykey")
	if !found {
		t.Error("expected to find key in cache")
	} else if val != "myval" {
		t.Errorf("expected 'myval', got %s", val)
	}
}

func TestCache_Delete(t *testing.T) {
	cache := jankdb.NewCache[int](5*time.Minute, 1*time.Minute)
	cache.Set("count", 42)
	_, found := cache.Get("count")
	if !found {
		t.Error("expected key to exist")
	}

	cache.Delete("count")
	_, found = cache.Get("count")
	if found {
		t.Error("expected key to be deleted")
	}
}
