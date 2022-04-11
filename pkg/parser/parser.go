package parser

import (
	"strings"
)

var (
	htmlEntities = map[string]string{"&lt;": "<", "&gt;": ">"}
)

func Lex(body string) string {
	var inTag bool
	var inBody bool
	var inEntity bool
	var tagName string
	var entityName string

	var text string
	for _, r := range body {
		c := string(r)

		switch true {
		// tags
		case c == "<":
			inTag = true
			if tagName == "/body" {
				inBody = false
			}
		case c == ">":
			inTag = false
			if strings.Contains(tagName, "body") {
				inBody = true
			}
			tagName = ""
		case inTag:
			tagName += c
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
		// body
		case inBody:
			text += c
		}
	}

	return text
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
