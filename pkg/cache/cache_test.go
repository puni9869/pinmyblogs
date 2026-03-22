package cache

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	c := NewCache[string, string]()
	c.Set("key1", "value1", 5*time.Second)

	got, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected key1 to exist")
	}
	if got != "value1" {
		t.Errorf("Get(key1) = %q, want %q", got, "value1")
	}
}

func TestGetMissing(t *testing.T) {
	c := NewCache[string, int]()

	_, ok := c.Get("missing")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestExpiry(t *testing.T) {
	c := NewCache[string, string]()
	c.Set("expire", "soon", 1*time.Millisecond)

	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("expire")
	if ok {
		t.Error("expected expired key to be gone")
	}
}

func TestDelete(t *testing.T) {
	c := NewCache[string, string]()
	c.Set("del", "me", 5*time.Second)
	c.Delete("del")

	_, ok := c.Get("del")
	if ok {
		t.Error("expected deleted key to be gone")
	}
}

func TestClear(t *testing.T) {
	c := NewCache[string, int]()
	c.Set("a", 1, 5*time.Second)
	c.Set("b", 2, 5*time.Second)
	c.Clear()

	if _, ok := c.Get("a"); ok {
		t.Error("expected cleared cache to have no key 'a'")
	}
	if _, ok := c.Get("b"); ok {
		t.Error("expected cleared cache to have no key 'b'")
	}
}

func TestOverwrite(t *testing.T) {
	c := NewCache[string, string]()
	c.Set("k", "v1", 5*time.Second)
	c.Set("k", "v2", 5*time.Second)

	got, ok := c.Get("k")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if got != "v2" {
		t.Errorf("Get(k) = %q, want %q", got, "v2")
	}
}

func TestIntKeys(t *testing.T) {
	c := NewCache[int, string]()
	c.Set(42, "answer", 5*time.Second)

	got, ok := c.Get(42)
	if !ok || got != "answer" {
		t.Errorf("Get(42) = %q, %v; want %q, true", got, ok, "answer")
	}
}
