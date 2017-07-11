package parse

import (
	"bufio"
	"bytes"
	"io"
)

type Token int

const (
	//Special
	ILLEGAL_TOKEN = iota
	EOF
	WS

	//Literals
	IDENT //todos and dates

	//Key Symbols
	STATUS_OPEN  //[
	STATUS_CLOSE //]

	//Misc
	ASTERISK //*
	COMMA    //,
	DOT      //.
	HASHTAG  //#
)

var eof = rune(0)

//Scanner represents a lexical scanner
type Scanner struct {
	r *bufio.Reader
}

//NewScanner returns a new instance of Scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

//read reads the next rune from the buffered reader.
//Returns the rune(0) if an error occurs(or io.EOF is returned).
func (s *Scanner) read() rune {
	r, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}

	return r
}

//unread places the previously read rune back to the reader
func (s *Scanner) unread() {
	s.r.UnreadRune()
}

//Scan returns the next token and its value
func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) || isDigit(ch) {
		s.unread()
		return s.scanIdent()
	}

	//Otherwise read individual character
	switch ch {
	case '#':
		return HASHTAG, "#"
	case '[':
		return STATUS_OPEN, "["
	case ']':
		return STATUS_CLOSE, "]"
	case ',':
		return COMMA, ","
	case '.':
		return DOT, "."
	case '*':
		return ASTERISK, "*"
	case eof:
		return EOF, string(ch)
	}

	return ILLEGAL_TOKEN, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	//Create buffer and read the current character into it
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	//Read every subsequent whitespace into the buffer.
	//Non Whitspace Characters and EOF will cause the loop to exit
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	//Read every subsequent ident character into the buffer.
	//Non ident Characters and EOF will cause the loop to exit

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
