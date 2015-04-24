package utils

import (
	"errors"
)

type State int

const (
	INITIAL State = iota
	NORMALCHAR
	SEPARATOR
	QUOTE
)

func Split(source string, separator rune, quote rune) ([]string, error) {
	var result []string
	var err error

	state := INITIAL
	var word []rune
	for _, c := range source {
		switch state {
		case INITIAL:
			if c == quote {
				state = QUOTE
				word = append(word, c)
			} else if c == separator {
				state = SEPARATOR
			} else {
				state = NORMALCHAR
				word = append(word, c)
			}
		case NORMALCHAR:
			if c == quote {
				state = QUOTE
				word = append(word, c)
			} else if c == separator {
				state = SEPARATOR
				result = append(result, string(word))
				word = []rune{}
			} else {
				word = append(word, c)
			}
		case QUOTE:
			word = append(word, c)
			if c == quote {
				state = NORMALCHAR
			}
		case SEPARATOR:
			if c == quote {
				state = QUOTE
				word = append(word, c)
			} else if c != separator {
				state = NORMALCHAR
				word = append(word, c)
			}
		default:
			state = INITIAL
			err = errors.New("PARSE ERROR : Unknown State.")
		}
	}

	if len(word) > 0 {
		result = append(result, string(word))
	}

	if state == QUOTE {
		err = errors.New("PARSE ERROR : Unmatched " + string(quote) + ".")
	}

	return result, err
}
