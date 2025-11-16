package hash_table

import (
	"testing"
	"time"
)

func TestDictSetGet(t *testing.T) {
	d := CreateDict()

	// --- Test Set & Get ---
	obj := d.NewObj("foo", "bar", 0)
	d.Set("foo", obj)

	v := d.Get("foo")
	if v == nil || v.Value != "bar" {
		t.Errorf("expected bar, got %v", v)
	}
}

func TestDictExpiry(t *testing.T) {
	d := CreateDict()

	// --- Set with TTL 50ms ---
	obj := d.NewObj("foo", "bar", 50)
	d.Set("foo", obj)

	// Should exist immediately
	if d.Get("foo") == nil {
		t.Errorf("expected value before expiry")
	}

	// Wait for expiry
	time.Sleep(60 * time.Millisecond)

	// Now should be expired
	if d.Get("foo") != nil {
		t.Errorf("expected nil after expiry")
	}
}

func TestDictDelete(t *testing.T) {
	d := CreateDict()

	obj := d.NewObj("foo", "bar", 0)
	d.Set("foo", obj)

	// Delete existing key
	deleted := d.Del("foo")
	if !deleted {
		t.Errorf("expected true when deleting existing key")
	}

	// Delete non-existing key
	deleted = d.Del("baz")
	if deleted {
		t.Errorf("expected false when deleting non-existing key")
	}
}

func TestDictGetExpiry(t *testing.T) {
	d := CreateDict()

	ttlMs := int64(1000)
	obj := d.NewObj("foo", "bar", ttlMs)
	d.Set("foo", obj)

	exp, exist := d.GetExpiry("foo")
	if !exist {
		t.Errorf("expected expiry to exist")
	}
	if exp <= uint64(time.Now().UnixMilli()) {
		t.Errorf("expected expiry in the future, got %v", exp)
	}
}
