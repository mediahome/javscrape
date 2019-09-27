package net

import "testing"

// TestNewCache ...
func TestNewCache(t *testing.T) {
	cache := NewCache("./tmp")
	t.Log(cache.Get("https://pics.javbus.com/cover/6qx9_b.jpg"))
}
