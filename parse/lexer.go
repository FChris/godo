package parse

import (
	"bufio"
	"bytes"
	"io"
)

//Token identifies the type of data that was read
type Token int

const (

	//illegal represents anything that cannot be identified by any other token
	illegal = iota

	//eof represents the end of file token
	eof

	//ws identifies a whitespace
	ws

	//ident are identifiers of todos and dates
	ident

	//statusOpen is the token for [
	statusOpen
	//statusClose is is the token for ]
	statusClose

	slash         // /
	semicolon     // ;
	colon         // :
	asterisk      // *
	comma         // ,
	dot           // .
	hashtag       // #
	bracket       // ( )
	currencySign  // $ €
	paragraph     // §
	ampersand     // &
	equals        // =
	tilde         // ~
	at            // @
	percent       // %
	dash          // -
	underscore    // _
)

var endoffile = rune(0)

//Scanner represents a lexical scanner
type scanner struct {
	*bufio.Reader
}

//NewScanner returns a new instance of Scanner
func NewScanner(r io.Reader) *scanner {
	return &scanner{bufio.NewReader(r)}
}

//read reads the next rune from the buffered reader.
//Returns the rune(0) if an error occurs(or io.eof is returned).
func (s *scanner) read() rune {
	r, _, err := s.ReadRune()
	if err != nil {
		return endoffile
	}

	return r
}

//Scan returns the next token and its value
func (s *scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.UnreadRune()
		return s.scanWhitespace()
	} else if isLetter(ch) || isDigit(ch) {
		s.UnreadRune()
		return s.scanIdent()
	}

	//Otherwise read individual character
	switch ch {
	case '#':
		return hashtag, "#"
	case '[':
		return statusOpen, "["
	case ']':
		return statusClose, "]"
	case ',':
		return comma, ","
	case '.':
		return dot, "."
	case ':':
		return colon, ":"
	case ';':
		return semicolon, ";"
	case '/':
		return slash, "/"
	case '*':
		return asterisk, "*"
	case '(':
		fallthrough
	case ')':
		return bracket, string(ch)
	case '~':
		return tilde, "~"
	case '€':
		fallthrough
	case '$':
		fallthrough
	case '£':
		fallthrough
	case '¥':
		return currencySign, string(ch)
	case '§':
		return paragraph, "§"
	case '&':
		return ampersand, "&"
	case '=':
		return equals, "="
	case '@':
		return at, "@"
	case '%':
		return percent, "%"
	case '-':
		return dash, "-"
	case '_':
		return underscore, "_"
	case endoffile:
		return eof, string(ch)
	}

	return illegal, string(ch)
}

func (s *scanner) scanWhitespace() (tok Token, lit string) {
	//Create buffer and read the current character into it
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	//Read every subsequent whitespace into the buffer.
	//Non Whitspace Characters and eof will cause the loop to exit
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.UnreadRune()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return ws, buf.String()
}

func (s *scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	//Read every subsequent ident character into the buffer.
	//Non ident Characters and eof will cause the loop to exit

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.UnreadRune()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return ident, buf.String()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == 'ä' || ch == 'Ö' ||
		ch == 'ö' || ch == 'Ä' ||
		ch == 'ü' || ch == 'Ü' ||
		ch == 'ß'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
