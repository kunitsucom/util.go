package stringz

import (
	"strings"
)

type _Case struct {
	allUpper      bool // [A-Z]+
	allLower      bool // [a-z]+
	firstUpper    bool // [A-Z][0-9a-z_-]+
	firstLower    bool // [a-z][0-9a-z_-]+
	hasUpper      bool // .*[A-Z].*
	hasLower      bool // .*[a-z].*
	hasUnderscore bool // .*_.*
	hasDash       bool // .*-.*
	hasNumber     bool // .*[0-9].*
	hasOther      bool // .*[^0-9a-zA-Z_-].*
}

func ParseCase(s string) _Case { //nolint:cyclop,revive
	c := _Case{
		allUpper: true,
		allLower: true,
	}

	for idx, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			c.allLower = false
			c.hasUpper = true
			if idx == 0 {
				c.firstUpper = true
			}
		case r >= 'a' && r <= 'z':
			c.allUpper = false
			c.hasLower = true
			if idx == 0 {
				c.firstLower = true
			}
		case r == '_':
			c.allUpper = false
			c.allLower = false
			c.hasUnderscore = true
		case r == '-':
			c.allUpper = false
			c.allLower = false
			c.hasDash = true
		case r >= '0' && r <= '9':
			c.allUpper = false
			c.allLower = false
			c.hasNumber = true
		default:
			c.allUpper = false
			c.allLower = false
			c.hasOther = true
		}
	}
	return c
}

// ex. snake_case, go, type_script, postgre_sql.
func IsSnakeCase(s string) bool {
	return ParseCase(s).IsSnakeCase()
}

func (c _Case) IsSnakeCase() bool {
	return c.hasUnderscore && c.hasLower && !c.hasUpper && !c.hasDash && !c.allLower
}

// ex. kebab-case, go, type-script, postgre-sql.
func IsKebabCase(s string) bool {
	return ParseCase(s).IsKebabCase()
}

func (c _Case) IsKebabCase() bool {
	return c.hasDash && c.hasLower && !c.hasUpper && !c.IsSnakeCase()
}

// ex. PascalCase, Go, TypeScript, PostgreSql.
func IsPascalCase(s string) bool {
	return ParseCase(s).IsPascalCase()
}

func (c _Case) IsPascalCase() bool {
	return c.firstUpper && !c.hasDash && !c.IsSnakeCase() && !c.IsKebabCase()
}

// ex. camelCase, go, typeScript, postgreSql.
func IsCamelCase(s string) bool {
	return ParseCase(s).IsCamelCase()
}

func (c _Case) IsCamelCase() bool {
	return c.firstLower && !c.hasDash && !c.IsSnakeCase() && !c.IsKebabCase() && !c.IsPascalCase()
}

func SplitSnakeCase(s string) []string {
	return strings.Split(s, "_")
}

func SplitKebabCase(s string) []string {
	return strings.Split(s, "-")
}

func SplitPascalCase(s string) []string {
	c := ParseCase(s)
	return c.splitPascalCase(s)
}

func (c _Case) splitPascalCase(s string) []string {
	return splitCamels(s)
}

func SplitCamelCase(s string) []string {
	c := ParseCase(s)
	return c.splitCamelCase(s)
}

func (c _Case) splitCamelCase(s string) []string {
	return splitCamels(s)
}

func SplitCase(s string) []string {
	c := ParseCase(s)

	switch {
	case c.IsSnakeCase():
		return SplitSnakeCase(s)
	case c.IsKebabCase():
		return SplitKebabCase(s)
	case c.IsPascalCase():
		return c.splitPascalCase(s)
	case c.IsCamelCase():
		return c.splitCamelCase(s)
	default:
		// as default, use camelCase.
		return c.splitCamelCase(s)
	}
}

func splitCamels(s string) []string { //nolint:cyclop
	var _words []string
	var _word []rune //nolint:prealloc

	var before rune
	for _, current := range s {
		if !('A' <= before && before <= 'Z' || before == 0) && 'A' <= current && current <= 'Z' {
			// aaA -> [aa A]
			//  ~~
			//	| \
			//	|  +-- current
			//  +-- before
			_words = append(_words, string(_word))
			_word = _word[:0]
		} else if 'A' <= before && before <= 'Z' && !('A' <= current && current <= 'Z') && len(_word) > 1 {
			// aAa -> [a Aa]
			//  ~~
			//	| \
			//	|  +-- current
			//  +-- before
			_words = append(_words, string(_word[:len(_word)-1]))
			_word = _word[len(_word)-1:]
		}
		_word = append(_word, current)
		before = current
	}
	if len(_word) > 0 {
		_words = append(_words, string(_word))
	}
	return _words
}
