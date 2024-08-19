package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func lexerHelper(lex *Lexer) ([]Token, []string) {
	var tokens []Token
	var literals []string
	for {
		tok := lex.Scan()
		tokens = append(tokens, tok.Token)
		literals = append(literals, tok.Literal)
		if tok.Token == EOF {
			return tokens, literals
		}
	}
}
func TestLexer(t *testing.T) {

	t.Run("scan into tokens succeeds", func(t *testing.T) {
		s := "MATCH (a:Person) RETURN a"
		lexer := NewLexerFromString(s)
		tokens, literals := lexerHelper(lexer)
		assert.Equal(t, []Token{MATCH, WS, OPEN_PARENTHESIS, STRING, DOUBLECOLON, STRING, CLOSED_PARENTHESIS, WS, RETURN, WS, STRING, EOF}, tokens)
		assert.Equal(t, []string{"MATCH", "", "(", "a", ":", "Person", ")", "", "RETURN", "", "a", ""}, literals)
	})

	t.Run("scan into tokens succeeds with quote", func(t *testing.T) {
		s := "Person{b:'c'}"
		lexer := NewLexerFromString(s)
		tokens, literals := lexerHelper(lexer)
		assert.Equal(t, []Token{STRING, OPEN_CURLYBRACKET, STRING, DOUBLECOLON, STRING, CLOSED_CURLYBRACKET, EOF}, tokens)
		assert.Equal(t, []string{"Person", "{", "b", ":", "c", "}", ""}, literals)
	})

	t.Run("scan into tokens failing with quote", func(t *testing.T) {
		s := "Person{b:'c}"
		lexer := NewLexerFromString(s)
		tokens, literals := lexerHelper(lexer)
		assert.Equal(t, []Token{STRING, OPEN_CURLYBRACKET, STRING, DOUBLECOLON, STRING, EOF}, tokens)
		assert.Equal(t, []string{"Person", "{", "b", ":", "c}", ""}, literals)
	})
}
