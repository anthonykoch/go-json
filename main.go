package main

import (
	"regexp"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	// "testing"
)

var S_SPACE = "[\u0009\u000a\u000d\u0020]"
var S_SIGN= "(?:+|-)"
var S_ONENINE = "[1-9]"
var S_DIGIT = "[0-9]"
var S_DIGITS = fmt.Sprintf("(?:%s+)", S_DIGIT)
var S_FRACTION = fmt.Sprintf("(?:\\.%s)", S_DIGITS)
var S_EXP = "(?:e|E)"
var S_INT = fmt.Sprintf("-?(?:%s%s|%s)", S_ONENINE, S_DIGITS, S_DIGIT)
var S_HEX = fmt.Sprintf("%s[a-fA-F]", S_DIGIT)
var S_ESCAPE = fmt.Sprintf("(?:[\"\\\\/bnrt]|u%s{4})", S_HEX)
var S_CHAR = fmt.Sprintf("(?:\\\\%s|[^\"\\\\])", S_ESCAPE)

var S_STRING = fmt.Sprintf("^\"%s*\"", S_CHAR)
var S_NULL = "^null"
var S_BOOL = "^(?:true|false)"
var S_NUMBER = fmt.Sprintf("^%s?%s?", S_INT, S_FRACTION)
var S_PUNCTUATOR = "^[\\[\\{\\}\\]:,]"

var R_SPACE = regexp.MustCompile(S_SPACE)
var R_STRING = regexp.MustCompile(S_STRING)
var R_NULL = regexp.MustCompile(S_NULL)
var R_BOOL = regexp.MustCompile(S_BOOL)
var R_NUMBER = regexp.MustCompile(S_NUMBER)
var R_PUNC = regexp.MustCompile(S_PUNCTUATOR)

var TYPE_NUMBER = "number"
var TYPE_STRING = "string"
var TYPE_BOOL = "boolean"
var TYPE_NULL = "null"
var TYPE_PUNC = "punctuator"

var tokenTypes = []string{
	TYPE_NUMBER,
	TYPE_STRING,
	TYPE_BOOL,
	TYPE_NULL,
	TYPE_PUNC,
}

var regexes = []*regexp.Regexp{
	R_NUMBER,
	R_STRING,
	R_BOOL,
	R_NULL,
	R_PUNC,
}

var log = fmt.Printf
var print = fmt.Println

func main() {
	spew.Dump(parse("{\"lul\": [1, \"hello\"]}"))
	spew.Dump(parse("123"))
	spew.Dump(parse("123.123"))
	spew.Dump(parse("\"hey\""))
	
	spew.Dump(parse("123 ]"))
}





type AST struct {
	program interface{}
}





func CreateBooleanLiteral(token Token) BooleanLiteral {
	var value bool 

	if (token.value == "true") {
		value = true		
	} else {
		value = false
	}

	return BooleanLiteral{
		name: "BooleanLiteral",
		value: value,
		start: token.start,
		end: token.end,
	}
}

type BooleanLiteral struct {
	name string
	value bool
	start int
	end int
}

type NullLiteral struct {
	name string
	start int
	end int
}

type StringLiteral struct {
	name string
	value string
	raw string
	start int
	end int
}

type NumericLiteral struct {
	name string
	value string
	raw string
	start int
	end int
}

type Property struct {
	name string
	value interface{}
	key string
}

type ObjectLiteral struct {
	name string
	properties []Property
	start int
	end int
}

type ArrayLiteral struct {
	name string
	elements []interface{}
	start int
	end int
}





type ParserError struct {
	message string
}

func (e ParserError) Error() string {
	return e.message
}

type Parser struct {
	lexer *JsonLexer
	source string
}

func CreateParser(input string) *Parser {
	parser := &Parser{ source: input }
	parser.lexer = CreateLexer(parser.source)
	
	return parser
}

func CreateNullLiteral(token Token) NullLiteral {
	return NullLiteral{
		name: "NullLiteral",
		start: token.start,
		end: token.end,
	}
}

func CreateStringLiteral(token Token) StringLiteral {
	return StringLiteral{
		name: "StringLiteral",
		// I'm just going to cheat here and not do this
		value: token.value,
		raw: token.value,
		start: token.start,
		end: token.end,
	}
}

func CreateNumericLiteral(token Token) NumericLiteral {
	return NumericLiteral{
		name: "NumericLiteral",
		// Going to cheat here as well
		value: token.value,
		raw: token.value,
		start: token.start,
		end: token.end,
	}
}





	
func (parser *Parser) ensure(value string) Token {
	// print("memed")
	// print(parser.lexer.peek())
	token, err := parser.lexer.next()
	// print(parser.lexer.peek())

	if err != nil { 
		panic(err)
	} else if token.value != value {
		panic(fmt.Sprintf("Token value '%s' does not match %s", token.value, value))
	} 

	return token
}


func (parser *Parser) ensureType(value string) Token {
	token, err := parser.lexer.next()

	if err != nil { 
		panic(err)
	} else if token.name != value {
		panic("Token value does not match")
	}

	return token
}


func (parser *Parser) is(value string) bool {
	token, err := parser.lexer.peek()

	return err != nil || token.value == value
}


func (parser *Parser) parseValue() interface{} {
	token, err := parser.lexer.peek()

	if err != nil {
		panic(err)
	}
	
	switch token.name {
		case TYPE_PUNC:
			if token.value == "[" {
				return parser.parseArrayLiteral()
			} else if token.value == "{" {
				return parser.parseObjectLiteral()
			} else {
				panic(fmt.Sprintf("Unexpected punctuator \"%s\"", token.value))
			}
		
		case TYPE_BOOL: 
			parser.lexer.next()
			return CreateBooleanLiteral(token)
		
		case TYPE_NULL:
			parser.lexer.next()
			return CreateNullLiteral(token)
		
		case TYPE_STRING:
			parser.lexer.next()
			return CreateStringLiteral(token)

		case TYPE_NUMBER:
			parser.lexer.next()
			return CreateNumericLiteral(token)
	}

	panic("Could not find a meme")
}


func (parser *Parser) parseArrayLiteral() ArrayLiteral {
	elements := make([]interface{}, 0)
	startToken := parser.ensure("[")
	// print(parser.lexer.peek())

	// print(parser.lexer.peek())
	
	for !parser.is("]") {
		node := parser.parseValue()
		elements = append(elements, node)

		if !parser.is("]") {
			parser.ensure(",")
		}
	}
	
	endToken := parser.ensure("]")

	return ArrayLiteral{
		name: "ArrayLiteral",
		elements: elements,
		start: startToken.start,
		end: endToken.end,
	}
}


func (parser *Parser) parseObjectLiteral() ObjectLiteral {
	properties := make([]Property, 0)

	startToken := parser.ensure("{")
	
	for !parser.is("}") {
		key := parser.ensureType(TYPE_STRING).value
		parser.ensure(":")
		property := Property{ key: key, value: parser.parseValue() }
		properties = append(properties, property)
	}
	
	endToken := parser.ensure("}")

	return ObjectLiteral{
		name: "ObjectLiteral", 
		properties: properties,
		start: startToken.start,
		end: endToken.end,
	}
}


func parse(input string) AST {
	parser := CreateParser(input)
	node := parser.parseValue()
	
	parser.lexer.peek()

	if len(parser.lexer.stash) > 0 {
		panic("Invalid JSON")
	}

	return AST{
		program: node,
	}
}





type LexError struct {
	message string
}

func (e LexError) Error() string {
	return e.message
}





type JsonLexer struct {
	input string
	source string
	index int
	stash []Token
}

func CreateLexer(input string) *JsonLexer {
	index := 0

	return &JsonLexer{
		stash: make([]Token, 0),
		input: input,
		source: input,
		index: index,
	}
}


func (lexer JsonLexer) hasEnded() bool {
	return lexer.index == len(lexer.source)
}


func (lexer *JsonLexer) lookahead(index int) (Token, error) {
	remainder := index - len(lexer.stash)
	
	for ; !lexer.hasEnded() && remainder > -1; remainder-- {
		token, err := lexer.lex()
		// print(lexer.hasEnded(), remainder, token)

		if err == nil {
			lexer.stash = append(lexer.stash, token)
		}
	}
	
	if index <= len(lexer.stash) - 1 {
		return lexer.stash[index], nil
	}

	return Token{}, LexError{ message: "lookahead index out of range" }
}


func (lexer *JsonLexer) peek() (Token, error) {
	return lexer.lookahead(0)
}


func (lexer *JsonLexer) next() (Token, error) {
	var token Token

	
	if len(lexer.stash) > 0 {
		token, lexer.stash = lexer.stash[0], lexer.stash[1:]
		// print("stash", token)
		return token, nil
	} else if (lexer.hasEnded()) {
		return token, &LexError{ message: "Lexer has already ended" }
	}

	token, err := lexer.lex()
	// print("lex", token)
	
	return token, err
}


func (lexer *JsonLexer) forward(times int) {
	start := lexer.index
	lexer.index += times
	lexer.input = lexer.source[start + times:len(lexer.source)]
}


func (lexer JsonLexer) getToken(re *regexp.Regexp, name string) Token {
	match := re.FindString(lexer.input)
	
	if match != "" {
		return CreateToken(name, match, lexer.index)
	}
	
	return Token{}
}


func (lexer *JsonLexer) skipWhitespace() {
	start := lexer.index
	i := start
	length := len(lexer.source)
	var char byte

	for ; i < length; {
		char = lexer.source[i]

		if char == '\u0020' || char == '\u0009' ||  char == '\u000a' || char == '\u000d' {
			
		} else {
			break
		}

		i++
	}

	lexer.forward(i - start)
}


// Will throw an "out of range" error if all tokens are consumed
func (lexer *JsonLexer) lex() (Token, error) {
	var token Token

	lexer.skipWhitespace()

	for i := 0; i < len(tokenTypes); i++ {
		if token := lexer.getToken(regexes[i], tokenTypes[i]); token.found {
			lexer.forward(len(token.value))
			lexer.skipWhitespace()

			// print(token, lexer.hasEnded(), lexer.index, len(lexer.source))

			return token, nil
		}
	}

	message := fmt.Sprintf("Unexpected token \"%s\"" , string(lexer.source[lexer.index]))
	
	return token, &LexError{ message: message }
}





type Token struct {
	found bool
	value string
	name string 
	start int 
	end int 
}

func CreateToken(name string, value string, start int) Token {
	return Token{
			name: name,
			value: value,
			found: true,
			start: start, 
			end: start + len(value) - 1,
	}
}





// func TestParser(t *testing.T) {
// 	total := Sum(5, 5)
// 	if total != 10 {
// 		 t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
// 	}
// }




