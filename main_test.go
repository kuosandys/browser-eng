package main

import (
	"bytes"
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
		buffer := bytes.Buffer{}
		show(&buffer, `<html><head></head><body><h1>hello bello</h1><img src="something.com" /><div><p>How are you</p></div></body></html>`)

		got := buffer.String()
		want := "hello belloHow are you"
		assertEqual(t, got, want)
	})

	t.Run("parse new line characters", func(t *testing.T) {
		buffer := bytes.Buffer{}
		show(&buffer, "<html><head></head><body><h1>line 1</h1><h2>\nline 2</h2></body></html>")

		got := buffer.String()
		want := "line 1\nline 2"
		assertEqual(t, got, want)
	})
}
