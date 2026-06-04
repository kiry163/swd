package core

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
	}

	got, mapping := normalizeWithMap("Ａ測@試", cfg)
	want := "a测试"
	wantMapping := []int{0, 1, 1, 1, 3, 3, 3}

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

func TestNormalizeWithMap_SimilarCharMapsCommonLeetspeak(t *testing.T) {
	got, mapping := normalizeWithMap("b@d t3st", Config{SimilarChar: true, IgnoreCase: true})
	want := "bad test"
	wantMapping := []int{0, 1, 2, 3, 4, 5, 6, 7}

	if got != want {
		t.Fatalf("normalizeWithMap() text = %q, want %q", got, want)
	}
	if !reflect.DeepEqual(mapping, wantMapping) {
		t.Fatalf("normalizeWithMap() mapping = %#v, want %#v", mapping, wantMapping)
	}
}

func TestNormalizeWithMap_IgnoreWidthFoldsHalfwidthKatakana(t *testing.T) {
	got, mapping := normalizeWithMap("ｶﾀｶﾅ", Config{IgnoreWidth: true})
	want := "カタカナ"
	wantMapping := []int{0, 0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3}

	if got != want {
		t.Fatalf("normalizeWithMap() text = %q, want %q", got, want)
	}
	if !reflect.DeepEqual(mapping, wantMapping) {
		t.Fatalf("normalizeWithMap() mapping = %#v, want %#v", mapping, wantMapping)
	}
}

func TestNormalizeWithMap_IgnoreWidthPreservesHalfwidthKatakanaMarkMapping(t *testing.T) {
	got, mapping := normalizeWithMap("ｶﾞ", Config{IgnoreWidth: true})
	want := "ガ"
	wantMapping := []int{0, 0, 0, 1, 1, 1}

	if got != want {
		t.Fatalf("normalizeWithMap() text = %q, want %q", got, want)
	}
	if !reflect.DeepEqual(mapping, wantMapping) {
		t.Fatalf("normalizeWithMap() mapping = %#v, want %#v", mapping, wantMapping)
	}
}
