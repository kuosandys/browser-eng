package parser

import (
	"reflect"
	"testing"
)

func TestShow(t *testing.T) {
	assertEqual := func(t testing.TB, got, want []interface{}) {
		t.Helper()
		if reflect.DeepEqual(got, want) == false {
			t.Errorf("got %q want %q", got, want)
		}
	}

	t.Run("parse some HTML", func(t *testing.T) {
		got := Lex(`<html><head></head><body><h1>hello bello</h1><div><p>How are you</p></div></body></html>`)
		want := []interface{}{
			&Tag{Tag: "html"},
			&Tag{Tag: "head"},
			&Tag{Tag: "/head"},
			&Tag{Tag: "body"},
			&Tag{Tag: "h1"},
			&Text{Text: "hello bello"},
			&Tag{Tag: "/h1"},
			&Tag{Tag: "div"},
			&Tag{Tag: "p"},
			&Text{Text: "How are you"},
			&Tag{Tag: "/p"},
			&Tag{Tag: "/div"},
			&Tag{Tag: "/body"},
			&Tag{Tag: "/html"},
		}

		assertEqual(t, got, want)
	})

	t.Run("parse new line characters", func(t *testing.T) {
		got := Lex("<html><head></head><body><h1>line 1</h1><h2>\nline 2</h2></body></html>")
		want := []interface{}{
			&Tag{Tag: "html"},
			&Tag{Tag: "head"},
			&Tag{Tag: "/head"},
			&Tag{Tag: "body"},
			&Tag{Tag: "h1"},
			&Text{Text: "line 1"},
			&Tag{Tag: "/h1"},
			&Tag{Tag: "h2"},
			&Text{Text: "\nline 2"},
			&Tag{Tag: "/h2"},
			&Tag{Tag: "/body"},
			&Tag{Tag: "/html"},
		}

		assertEqual(t, got, want)
	})

	t.Run("output only what's in the body Tag", func(t *testing.T) {
		got := Lex("<html><head><title>some title</title></head><body><h1>hello from the body</h1></body></html>")
		want := []interface{}{
			&Tag{Tag: "html"},
			&Tag{Tag: "head"},
			&Tag{Tag: "title"},
			&Tag{Tag: "/title"},
			&Tag{Tag: "/head"},
			&Tag{Tag: "body"},
			&Tag{Tag: "h1"},
			&Text{Text: "hello from the body"},
			&Tag{Tag: "/h1"},
			&Tag{Tag: "/body"},
			&Tag{Tag: "/html"},
		}

		assertEqual(t, got, want)
	})

	t.Run("parse the less-than entity", func(t *testing.T) {
		got := Lex("<html><head></head><body>&lt;3</body></html>")
		want := []interface{}{
			&Tag{Tag: "html"},
			&Tag{Tag: "head"},
			&Tag{Tag: "/head"},
			&Tag{Tag: "body"},
			&Text{Text: "<3"},
			&Tag{Tag: "/body"},
			&Tag{Tag: "/html"},
		}

		assertEqual(t, got, want)
	})

	t.Run("parse the greater-than entity", func(t *testing.T) {
		got := Lex("<html><head></head><body>--&gt;</body></html>")
		want := []interface{}{
			&Tag{Tag: "html"},
			&Tag{Tag: "head"},
			&Tag{Tag: "/head"},
			&Tag{Tag: "body"},
			&Text{Text: "-->"},
			&Tag{Tag: "/body"},
			&Tag{Tag: "/html"},
		}

		assertEqual(t, got, want)
	})

	t.Run("handles unclosed tags", func(t *testing.T) {
		got := Lex("<html><head></head><body>Hi!<hr</body></html>")
		want := []interface{}{
			&Tag{Tag: "html"},
			&Tag{Tag: "head"},
			&Tag{Tag: "/head"},
			&Tag{Tag: "body"},
			&Text{Text: "Hi!"},
			&Tag{Tag: "/body"},
			&Tag{Tag: "/html"},
		}

		assertEqual(t, got, want)
	})

	t.Run("handles empty tags", func(t *testing.T) {
		got := Lex("<html><head></head><body><></></body></html>")
		want := []interface{}{
			&Tag{Tag: "html"},
			&Tag{Tag: "head"},
			&Tag{Tag: "/head"},
			&Tag{Tag: "body"},
			&Tag{Tag: ""},
			&Tag{Tag: "/"},
			&Tag{Tag: "/body"},
			&Tag{Tag: "/html"},
		}

		assertEqual(t, got, want)
	})

	t.Run("handles plain strings", func(t *testing.T) {
		got := Lex("boop")
		want := []interface{}{
			&Text{Text: "boop"},
		}

		assertEqual(t, got, want)
	})
}
