package engine

import (
	"reflect"
	"testing"
)

func TestNormalizeWithMap_AppliesEnabledOptionsAndPreservesRuneMapping(t *testing.T) {
	cfg := Config{
		IgnoreSymbol: true,
		IgnoreWidth:  true,
		IgnoreCase:   true,
		Traditional:  true,
		SimilarChar:  true,
		Pinyin:       true,
	}

	got, mapping := normalizeWithMap("Ａ測@試", cfg)
	want := "aceshi"
	wantMapping := []int{0, 1, 1, 3, 3, 3}

	if got != want {
		t.Fatalf("normalizeWithMap() text = %q, want %q", got, want)
	}
	if !reflect.DeepEqual(mapping, wantMapping) {
		t.Fatalf("normalizeWithMap() mapping = %#v, want %#v", mapping, wantMapping)
	}
}

func TestNormalizeWithMap_LeavesDisabledOptionsUntouched(t *testing.T) {
	got, mapping := normalizeWithMap("Ａ測@試", Config{})
	want := "Ａ測@試"
	wantMapping := []int{0, 0, 0, 1, 1, 1, 2, 3, 3, 3}

	if got != want {
		t.Fatalf("normalizeWithMap() text = %q, want %q", got, want)
	}
	if !reflect.DeepEqual(mapping, wantMapping) {
		t.Fatalf("normalizeWithMap() mapping = %#v, want %#v", mapping, wantMapping)
	}
}

func TestNormalizeWithMap_HomophoneUsesSimilarMapForCurrentBehavior(t *testing.T) {
	got, mapping := normalizeWithMap("零一", Config{Homophone: true})
	want := "01"
	wantMapping := []int{0, 1}

	if got != want {
		t.Fatalf("normalizeWithMap() text = %q, want %q", got, want)
	}
	if !reflect.DeepEqual(mapping, wantMapping) {
		t.Fatalf("normalizeWithMap() mapping = %#v, want %#v", mapping, wantMapping)
	}
}
