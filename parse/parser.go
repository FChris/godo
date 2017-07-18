package parse

import (
	"fmt"
	"io"
	"time"
	"bytes"
	"strings"
	"github.com/fchris/godo/task"
)

const Timeformat string = "02.01.06"

type Parser struct {
	s *scanner
	buf struct {
		tok Token  //last read token
		lit string //last string literal
		n   int    //buffer size (max = 1)
	}
}

//NewParser returns an instance of a new parser
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.s.Scan()
	if tok == WS {
		tok, lit = p.scanIgnoreWhitespace()
	}

	return
}

// Parse returns a the next task.Day struct that can be parsed from the input or an error if no new task.Day can be
// parsed
func (p *Parser) Parse() (task.Day, error) {
	var	taskDay task.Day

	tok, lit := p.scanIgnoreWhitespace()
	if tok != HASHTAG && tok != EOF {
		return taskDay, fmt.Errorf("found %q, expected #", lit)
	}

	if tok == EOF {
		return taskDay, nil
	}

	var buf bytes.Buffer
	for {
		//Read a field
		tok, lit := p.s.Scan()

		if tok != IDENT && tok != DOT && tok != WS {
			return taskDay, fmt.Errorf("found %q, expected field or dot", lit)
		}

		if tok == WS && buf.Len() > 0 {
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
		if tok == HASHTAG || tok == EOF {
			return taskDay, nil
		}

		if tok != STATUS_OPEN {
			return taskDay, fmt.Errorf("found %q, expected [", lit)
		}

		if tok, lit := p.s.Scan(); tok != WS && tok != IDENT {
			return taskDay, fmt.Errorf("found %q, expected WS or X", lit)
		} else {
			if tok == IDENT {
				todo.Complete = true
			} else {
				todo.Complete = false
			}
		}

		if tok, lit := p.s.Scan(); tok != STATUS_CLOSE {
			return taskDay, fmt.Errorf("found %q, expected ]", lit)
		}

		var buf bytes.Buffer

		for {
			//Read a field
			tok, lit := p.s.Scan()

			if  !isDescriptionToken(tok) && tok != STATUS_OPEN && tok != EOF && tok != HASHTAG {
				return taskDay, fmt.Errorf("found %q, expected field", lit)
			}

			if tok == EOF || tok == HASHTAG {
				todo.Description = strings.Trim(buf.String(), " \n")
				taskDay.Todos.InsertTodo(*todo)
				p.s.UnreadRune()
				return taskDay, nil
			}

			if tok == STATUS_OPEN {
				p.s.UnreadRune()
				break
			}

			buf.WriteString(lit)
		}

		todo.Description = strings.Trim(buf.String(), " \n")
		taskDay.Todos.InsertTodo(*todo)
	}
}

func isDescriptionToken(tok Token) bool {
	return tok == WS ||
		tok == IDENT ||
		tok == DOT ||
		tok == COMMA ||
		tok == SLASH ||
		tok == SEMICOLON ||
		tok == COLON ||
		tok == ASTERISK ||
		tok == BRACKET ||
		tok == CURRENCY_SIGN ||
		tok == PARAGRAPH ||
		tok == DASH ||
		tok == UNDERSCORE
}
