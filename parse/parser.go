package parse

import (
	"bytes"
	"fmt"
	"github.com/fchris/towg/task"
	"io"
	"strings"
	"time"
)

//Timeformat describes the format used to parse dates
const Timeformat string = "02.01.06"

//Parser provides the functionality to parse files that were tokenized by lexer
type Parser struct {
	*scanner
}

//NewParser returns an instance of a new parser
func NewParser(r io.Reader) *Parser {
	return &Parser{NewScanner(r)}
}

func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.Scan()
	if tok == ws {
		tok, lit = p.scanIgnoreWhitespace()
	}

	return
}

// Parse returns a the next task.Day struct that can be parsed from the input or an error if no new task.Day can be
// parsed
func (p *Parser) Parse() (task.Day, error) {
	var taskDay task.Day

	tok, lit := p.scanIgnoreWhitespace()
	if tok != hashtag && tok != eof {
		return taskDay, fmt.Errorf("found %q, expected #", lit)
	}

	if tok == eof {
		return taskDay, nil
	}

	var buf bytes.Buffer
	for {
		//Read a field
		tok, lit := p.Scan()

		if tok != ident && tok != dot && tok != ws {
			return taskDay, fmt.Errorf("found %q, expected field or dot", lit)
		}

		if tok == ws && buf.Len() > 0 {
			dateString := strings.Trim(buf.String(), " ")
			dueTime, err := time.Parse(Timeformat, dateString)
			if err != nil {
				return taskDay, err
			}
			taskDay.Date = dueTime
			break
		}

		buf.WriteString(lit)
	}
	buf.Reset()

	for {
		todo := &task.Todo{}

		tok, lit := p.scanIgnoreWhitespace()
		if tok == hashtag || tok == eof {
			return taskDay, nil
		}

		if tok != dash {
			return taskDay, fmt.Errorf("found %q, expected -", lit)
		}

		tok, lit = p.scanIgnoreWhitespace()
		if tok != statusOpen {
			return taskDay, fmt.Errorf("found %q, expected [", lit)
		}

		if tok, lit = p.Scan(); tok != ws && tok != ident {
			return taskDay, fmt.Errorf("found %q, expected ws or X", lit)
		}

		if tok == ident {
			todo.Complete = true
		} else {
			todo.Complete = false
		}

		if tok, lit := p.Scan(); tok != statusClose {
			return taskDay, fmt.Errorf("found %q, expected ]", lit)
		}

		var buf bytes.Buffer

		for {
			//Read a field
			tok, lit := p.Scan()

			if !isDescriptionToken(tok) && tok != dash && tok != eof && tok != hashtag {
				return taskDay, fmt.Errorf("found %q, expected field", lit)
			}

			if tok == eof || tok == hashtag {
				todo.Description = strings.Trim(buf.String(), " \n")
				taskDay.Todos.InsertTodo(*todo)
				p.UnreadRune()
				return taskDay, nil
			}

			if tok == dash {
				p.UnreadRune()
				break
			}

			buf.WriteString(lit)
		}

		todo.Description = strings.Trim(buf.String(), " \n")
		taskDay.Todos.InsertTodo(*todo)
	}
}

func isDescriptionToken(tok Token) bool {
	return tok == ws ||
		tok == ident ||
		tok == dot ||
		tok == comma ||
		tok == slash ||
		tok == semicolon ||
		tok == colon ||
		tok == asterisk ||
		tok == bracket ||
		tok == currencySign ||
		tok == paragraph ||
		tok == underscore
}
