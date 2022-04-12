package parser

import (
	"strings"
)

var (
	htmlEntities = map[string]string{"&lt;": "<", "&gt;": ">"}
)

type Text struct {
	Text string
}

// NewText creates a new Text object with the specified Text.
func NewText(text string) *Text {
	return &Text{Text: text}
}

type Tag struct {
	Tag string
}

// NewTag creates a new Tag object with the specified Tag.
func NewTag(tag string) *Tag {
	return &Tag{Tag: tag}
}

type Token interface {
	*Tag | *Text
}

func Lex(body string) []interface{} {
	var inTag bool
	var inBody bool
	var inEntity bool
	var entityName string
	var out []interface{}
	var text string
	for _, r := range body {
		c := string(r)

		switch true {
		// tags
		case c == "<":
			if len(text) > 0 && !inTag && inBody {
				out = append(out, NewText(text))
			}
			text = ""
			inTag = true
			if text == "/body" {
				inBody = false
			}
		case c == ">":
			inTag = false
			if strings.Contains(text, "body") {
				inBody = true
			}
			out = append(out, NewTag(text))
			text = ""
		case inTag:
			text += c
		// entities
		case c == "&":
			inEntity = true
			entityName += c
		case inEntity && c == ";":
			entityName += c
			character := htmlEntities[entityName]
			text += character
			inEntity = false
			entityName = ""
		case inBody && inEntity:
			entityName += c
		default:
			text += c
		}
	}

	// dump any accumulated Text
	if !inTag && len(text) > 0 {
		out = append(out, NewText(text))
	}

	return out
}

func Transform(body string) string {
	bodyTransformed := "<body>"

	for _, r := range body {
		c := string(r)
		switch true {
		case c == "<":
			bodyTransformed += "&lt;"
		case c == ">":
			bodyTransformed += "&gt;"
		default:
			bodyTransformed += c
		}
	}

	return bodyTransformed + "</body>"
}
