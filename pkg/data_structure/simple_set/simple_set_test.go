package simple_set

import (
	"reflect"
	"sort"
	"testing"
)

func TestSimpleSet(t *testing.T) {
	set := NewSimpleSet("test")

	// --- Test SADD ---
	added := set.Add("a", "b", "c", "a") // "a" bị duplicate
	if added != 3 {
		t.Errorf("expected 3 added, got %d", added)
	}

	// --- Test SISMEMBER ---
	if set.IsMember("a") != 1 {
		t.Errorf("expected a to be member")
	}
	if set.IsMember("x") != 0 {
		t.Errorf("expected x not to be member")
	}

	// --- Test SMEMBERS ---
	members := set.Members()
	sort.Strings(members) // để so sánh không phụ thuộc thứ tự
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(members, expected) {
		t.Errorf("expected %v, got %v", expected, members)
	}

	// --- Test SREM ---
	removed := set.Rem("a", "x") // "a" exist, "x" không exist
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}
	if set.IsMember("a") != 0 {
		t.Errorf("expected a to be removed")
	}
}
