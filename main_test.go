package main

import (
	"testing"
)

func TestShow(t *testing.T) {
	assertEqual := func(t testing.TB, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}

	t.Run("parse some HTML", func(t *testing.T) {
		got := show(`<html><head></head><body><h1>hello bello</h1><img src="something.com" /><div><p>How are you</p></div></body></html>`)
		want := "hello belloHow are you"
		assertEqual(t, got, want)
	})

	t.Run("parse new line characters", func(t *testing.T) {
		got := show("<html><head></head><body><h1>line 1</h1><h2>\nline 2</h2></body></html>")
		want := "line 1\nline 2"
		assertEqual(t, got, want)
	})

	t.Run("output only what's in the body tag", func(t *testing.T) {
		got := show("<html><head><title>some title</title></head><body><h1>hello from the body</h1></body></html>")
		want := "hello from the body"
		assertEqual(t, got, want)
	})

	t.Run("parse the less-than entity", func(t *testing.T) {
		got := show("<html><head></head><body>&lt;3</body></html>")
		want := "<3"
		assertEqual(t, got, want)
	})

	t.Run("parse the greater-than entity", func(t *testing.T) {
		got := show("<html><head></head><body>--&gt;</body></html>")
		want := "-->"
		assertEqual(t, got, want)
	})
}
