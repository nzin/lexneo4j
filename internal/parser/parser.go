package parser

import (
	"fmt"
	"strings"
)

// Parser represents a parser, including a scanner and the underlying raw input.
// It also contains a small buffer to allow for two unscans.
type Parser struct {
	s   *Lexer
	raw string
	buf TokenStack
}

// NewParser returns a new instance of Parser.
func NewParser(s string) *Parser {
	return &Parser{s: NewLexer(strings.NewReader(s)), raw: s}
}

// Parse takes the raw string and returns the root node of the AST.
func (p *Parser) Parse() (*CypherQuery, error) {
	operation, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	return operation, nil
}

// parseQuery parse stuff like MATCH (n:Person{foo:'bar'}) RETURN n.foo"
func (p *Parser) parseQuery() (*CypherQuery, error) {
	tok, _ := p.scanIgnoreWhitespace()
	if tok != MATCH {
		return nil, fmt.Errorf("not able to find a MATCH at the beginning of the expression")
	}

	node, err := p.parseNode()
	if err != nil {
		return nil, err
	}

	cypher := CypherQuery{MatchNode: *node}

	tok, _ = p.scanIgnoreWhitespace()

	// -->
	if tok == RELATIONSHIP {
		rel := CypherRelationShip{}
		cypher.Relationship = &rel

		var lit string
		tok, lit = p.scanIgnoreWhitespace()
		// relationship props to scan
		if tok != TO_RELATIONSHIP && tok != RELATIONSHIP {
			p.unscan(TokenInfo{Token: tok, Literal: lit})

			relProps, err := p.parseRelationshipProperties()
			if err != nil {
				return nil, err
			}
			rel.Props = relProps

			tok, lit = p.scanIgnoreWhitespace()
			if tok != TO_RELATIONSHIP && tok != RELATIONSHIP {
				return nil, fmt.Errorf("expected '->' or '-'. Got %s (%d)", lit, tok)
			}
			if tok == TO_RELATIONSHIP {
				rel.Direction = REL_TO
			} else {
				rel.Direction = REL_BOTH
			}
		}

		// and the target node
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		rel.Target = *node

		tok, _ = p.scanIgnoreWhitespace()
	}

	// <--
	if tok == FROM_RELATIONSHIP {
		rel := CypherRelationShip{
			Direction: REL_FROM,
		}
		cypher.Relationship = &rel

		var lit string
		tok, lit = p.scanIgnoreWhitespace()
		if tok != RELATIONSHIP {
			p.unscan(TokenInfo{Token: tok, Literal: lit})

			relProps, err := p.parseRelationshipProperties()
			if err != nil {
				return nil, err
			}
			rel.Props = relProps

			tok, _ = p.scanIgnoreWhitespace()
			if tok != RELATIONSHIP {
				return nil, fmt.Errorf("expected '-'")
			}
		}

		// and the target node
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		rel.Target = *node

		tok, _ = p.scanIgnoreWhitespace()
	}

	if tok == RETURN {
		ret, err := p.parseReturn()
		if err != nil {
			return nil, err
		}
		cypher.Return = ret
	}

	return &cypher, nil
}

// parseReturn scans stuff like "a,b.propname"
func (p *Parser) parseReturn() (CypherReturn, error) {
	ret := CypherReturn{}
	tok, lit := p.scanIgnoreWhitespace()

	for tok != EOF {
		if tok != STRING {
			return nil, fmt.Errorf("not able to find a correct return definition (return element name missing: %s)", lit)
		}
		elementName := lit
		retElement := CypherVariableReturn{
			VariableName: elementName,
		}

		tok, _ = p.scanIgnoreWhitespace()

		if tok == DOT {
			tok, lit = p.scanIgnoreWhitespace()
			if tok != STRING {
				return nil, fmt.Errorf("not able to find a correct return definition (return element property missing: %s)", lit)
			}
			elementProperty := lit
			retElement.Property = &elementProperty
			tok, _ = p.scanIgnoreWhitespace()
		}

		if tok != EOF && tok != COMMA {
			return nil, fmt.Errorf("not able to find a correct return definition (comma expected)")
		}
		ret = append(ret, retElement)

		if tok == COMMA {
			tok, lit = p.scanIgnoreWhitespace()
			if tok == EOF {
				return nil, fmt.Errorf("missing return value after comma)")
			}
		}
	}
	return ret, nil
}

// parseNode scans stuff like "(a:Person{foo:'bar'})"
func (p *Parser) parseNode() (*CypherNode, error) {

	node := CypherNode{}

	tok, _ := p.scanIgnoreWhitespace()
	if tok != OPEN_PARENTHESIS {
		return nil, fmt.Errorf("not able to find a correct node definition")
	}
	tok, lit := p.scanIgnoreWhitespace()

	if tok == CLOSED_PARENTHESIS {
		// empty "()"" node definition
		return &node, nil
	}

	if tok == STRING {
		variableName := lit
		node.VariableName = &variableName

		tok, lit = p.scanIgnoreWhitespace()
		if tok == CLOSED_PARENTHESIS {
			// end, i.e. "(n)"
			return &node, nil
		}
	}

	if tok == DOUBLECOLON {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != STRING {
			return nil, fmt.Errorf("missing type definition after ':'")
		}
		typeName := lit
		node.TypeName = &typeName

		tok, lit = p.scanIgnoreWhitespace()
		if tok == CLOSED_PARENTHESIS {
			// end, i.e. "(n:t)"
			return &node, nil
		}
	}

	if tok != OPEN_CURLYBRACKET {
		return nil, fmt.Errorf("missing '{ ... }'")
	}

	p.unscan(TokenInfo{Token: tok, Literal: lit})
	props, err := p.parseProperties()
	if err != nil {
		return nil, err
	}
	node.Props = props

	tok, _ = p.scanIgnoreWhitespace()
	if tok != CLOSED_PARENTHESIS {
		return nil, fmt.Errorf("not able to find a correct node definition")
	}
	return &node, nil
}

// parseNode scans stuff like "(a:Person{foo:'bar'})"
func (p *Parser) parseRelationshipProperties() (*CypherNode, error) {

	node := CypherNode{}

	tok, _ := p.scanIgnoreWhitespace()
	if tok != OPEN_BRACKET {
		return nil, fmt.Errorf("not able to find a correct relationship definition")
	}
	tok, lit := p.scanIgnoreWhitespace()

	if tok == CLOSED_BRACKET {
		// empty "()"" node definition
		return &node, nil
	}

	if tok == STRING {
		variableName := lit
		node.VariableName = &variableName

		tok, lit = p.scanIgnoreWhitespace()
		if tok == CLOSED_BRACKET {
			// end, i.e. "(n)"
			return &node, nil
		}
	}

	if tok == DOUBLECOLON {
		tok, lit = p.scanIgnoreWhitespace()
		if tok != STRING {
			return nil, fmt.Errorf("missing type definition after ':'")
		}
		typeName := lit
		node.TypeName = &typeName

		tok, lit = p.scanIgnoreWhitespace()
		if tok == CLOSED_BRACKET {
			// end, i.e. "(n:t)"
			return &node, nil
		}
	}

	if tok != OPEN_CURLYBRACKET {
		return nil, fmt.Errorf("missing '{ ... }'")
	}

	p.unscan(TokenInfo{Token: tok, Literal: lit})
	props, err := p.parseProperties()
	if err != nil {
		return nil, err
	}
	node.Props = props

	tok, _ = p.scanIgnoreWhitespace()
	if tok != CLOSED_BRACKET {
		return nil, fmt.Errorf("not able to find a correct node definition")
	}
	return &node, nil
}

// parseProperties scans stuff like "{foo:'bar'}"
func (p *Parser) parseProperties() (map[string]string, error) {
	props := make(map[string]string)

	tok, _ := p.scanIgnoreWhitespace()
	if tok != OPEN_CURLYBRACKET {
		return nil, fmt.Errorf("not able to find a correct properties definition")
	}
	tok, lit := p.scanIgnoreWhitespace()

	for tok != CLOSED_CURLYBRACKET && tok != EOF {
		if tok != STRING {
			return nil, fmt.Errorf("not able to find a correct properties definition (property name missng)")
		}
		propName := lit

		tok, _ = p.scanIgnoreWhitespace()
		if tok != DOUBLECOLON {
			return nil, fmt.Errorf("not able to find a correct properties definition (double colon missing)")
		}

		tok, lit = p.scanIgnoreWhitespace()
		if tok != STRING {
			return nil, fmt.Errorf("not able to find a correct properties definition (property value missing)")
		}
		props[propName] = lit

		tok, lit = p.scanIgnoreWhitespace()
		if tok != CLOSED_CURLYBRACKET && tok != COMMA {
			return nil, fmt.Errorf("not able to find a correct properties definition (comma or curly bracket missing)")
		}
	}

	return props, nil

}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.Len() != 0 {
		// Can ignore the error since it's not empty.
		tokenInf, _ := p.buf.Pop()
		return tokenInf.Token, tokenInf.Literal
	}

	// Otherwise read the next token from the scanner.
	tokenInf := p.s.Scan()
	tok, lit = tokenInf.Token, tokenInf.Literal
	return tok, lit
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return tok, lit
}

// unscan pushes the previously read tokens back onto the buffer.
func (p *Parser) unscan(tok TokenInfo) {
	p.buf.Push(tok)
}
