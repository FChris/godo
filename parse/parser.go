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
	s *Scanner
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

func (p *Parser) scan() (tok Token, lit string) {
	//if we have a token on the buffer, return it
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	//Otherwise read the next token from the scanner
	tok, lit = p.s.Scan()

	//Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan() //shouldn't we use scanIgnore whitespace again?
	}

	return
}

func (p *Parser) Parse() (*task.Day, error) {
	taskDay := &task.Day{} //What are we actually doing here?

	if tok, lit := p.scanIgnoreWhitespace(); tok != HASHTAG {
		return nil, fmt.Errorf("found %q, expected #", lit)
	}

	var buf bytes.Buffer
	for {
		//Read a field
		tok, lit := p.scan()

		if tok != IDENT && tok != DOT && tok != WS {
			return nil, fmt.Errorf("found %q, expected field or dot", lit)
		}

		if tok == WS && buf.Len() > 0 {
			dateString := strings.Trim(buf.String(), " ")
			fmt.Println(dateString)
			dueTime, err := time.Parse(Timeformat, dateString)
			if err != nil {
				return nil, err
			}
			taskDay.Date = dueTime
			break
		}

		buf.WriteString(lit)
	}
	buf.Reset()

	for {
		task := &task.Todo{}

		tok, lit := p.scanIgnoreWhitespace()
		if tok == HASHTAG || tok == EOF {
			return taskDay, nil
		}

		if tok != STATUS_OPEN {
			return nil, fmt.Errorf("found %q, expected [", lit)
		}

		if tok, lit := p.scan(); tok != WS && tok != IDENT {
			return nil, fmt.Errorf("found %q, expected WS or X", lit)
		} else {
			if tok == IDENT {
				task.Complete = true
			} else {
				task.Complete = false
			}
		}

		if tok, lit := p.scan(); tok != STATUS_CLOSE {
			return nil, fmt.Errorf("found %q, expected ]", lit)
		}

		var buf bytes.Buffer

		for {
			//Read a field
			tok, lit := p.scan()

			if tok != WS && tok != IDENT && tok != DOT && tok != COMMA && tok != STATUS_OPEN && tok != EOF {
				return nil, fmt.Errorf("found %q, expected field", lit)
			}

			if tok == EOF || tok == STATUS_OPEN {
				p.s.unread()
				break
			}

			buf.WriteString(lit)
		}

		task.Description = strings.Trim(buf.String(), " \n")
		taskDay.Todos = append(taskDay.Todos, *task)
	}
}
