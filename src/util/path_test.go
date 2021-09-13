package util

import "testing"

func TestHasUrlPrefixDir(t *testing.T) {
	var full, prefix string

	full = "/a/b/c"
	prefix = "/a/b"
	if !HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/"
	if !HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c"
	if !HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/e/"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/cd"
	prefix = "/a/b/c"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/c/"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/d"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}

	full = "/a/b/c"
	prefix = "/a/b/de"
	if HasUrlPrefixDir(full, prefix) {
		t.Error(full, prefix)
	}
}
