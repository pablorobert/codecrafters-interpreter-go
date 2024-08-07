package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Token struct {
	Pos   int
	Type  string
	Token string
	Value *string
}

var posToken int
var line int
var tokens []Token
var hasError bool

/*
and, class, else, false, for, fun,
if, nil, or, print, return, super, this, true, var, while
*/
var reserverd map[string]string = map[string]string{
	"and":    "AND",
	"class":  "CLASS",
	"else":   "ELSE",
	"false":  "FALSE",
	"for":    "FOR",
	"fun":    "FUN",
	"if":     "IF",
	"nil":    "NIL",
	"or":     "OR",
	"print":  "PRINT",
	"return": "RETURN",
	"super":  "SUPER",
	"this":   "THIS",
	"true":   "TRUE",
	"var":    "VAR",
	"while":  "WHILE",
}

func main() {
	posToken = -1
	line = 1
	hasError = false

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" && command != "parse" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	for posToken < len(fileContents)-1 {
		newToken, err, ignored := scanToken(string(fileContents))
		if ignored { //spaces, tabs, comments
			continue
		} else if err != nil {
			hasError = true
		} else {
			tokens = append(tokens, newToken)
		}
	}
	token := Token{posToken, "EOF", "", nil}
	tokens = append(tokens, token)

	if command == "tokenize" {
		for token := range tokens {
			printToken(tokens[token])
		}
	}
	if command == "parse" {
		for token := range tokens {
			printTokenValue(tokens[token])
		}
	}

	if hasError {
		os.Exit(65)
	}
}

func scanToken(input string) (Token, error, bool) {
	c, _ := advance(input)
	ignored := false
	switch {
	case isAlpha(c):
		word := c
		next := peek(input, 1)
		if next != "" && (isAlpha(next) || isDigit(next)) {
			word += eatIdentifier(input)
		}
		if val, ok := reserverd[word]; ok {
			return Token{posToken, val, word, nil}, nil, ignored
		}
		return Token{posToken, "IDENTIFIER", word, nil}, nil, ignored
	case isDigit(c):
		num := c
		var err error
		next := peek(input, 1)
		if next != "" && (isDigit(next) || next == ".") {
			if next == "." {
				num += next
				posToken++
				posToken++
			}
			num, err, _ = eatNumber(input) //including '.'
			if err != nil {
				//fmt.Fprintf(os.Stderr, "[line %d] Error: number.", line)
				return Token{}, errors.New("Error"), ignored
			}
			num = c + num
		}
		numValue := strings.Clone(num)
		if !strings.Contains(num, ".") {
			numValue += ".0"
		}
		numValue = strings.Replace(numValue, ".00", ".0", -1) //workaround
		return Token{posToken, "NUMBER", num, &numValue}, nil, ignored
	case c[0] == '"':
		str, err := eatString(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.", line)
			return Token{}, errors.New("Error"), ignored
		} else {
			return Token{posToken, "STRING", fmt.Sprintf("\"%s\"", str), &str}, nil, ignored
		}
	case c[0] == '\n':
		line++
		ignored = true
		return Token{}, nil, ignored
	case strings.TrimSpace(c) == "": //spaces, tab
		ignored = true
		return Token{}, nil, ignored
	case c == "=":
		token := Token{posToken, "EQUAL", "=", nil}
		next := peek(input, 1)
		if next == "=" {
			token = Token{posToken, "EQUAL_EQUAL", "==", nil}
			posToken++
		}
		return token, nil, ignored
	case c == "!":
		token := Token{posToken, "BANG", "!", nil}
		next := peek(input, 1)
		if next == "=" {
			token = Token{posToken, "BANG_EQUAL", "!=", nil}
			posToken++
		}
		return token, nil, ignored
	case c == "<":
		token := Token{posToken, "LESS", "<", nil}
		next := peek(input, 1)
		if next == "=" {
			token = Token{posToken, "LESS_EQUAL", "<=", nil}
			posToken++
		}
		return token, nil, ignored
	case c == ">":
		token := Token{posToken, "GREATER", ">", nil}
		next := peek(input, 1)
		if next == "=" {
			token = Token{posToken, "GREATER_EQUAL", ">=", nil}
			posToken++
		}
		return token, nil, ignored
	case c == "/":
		token := Token{posToken, "SLASH", "/", nil}
		next := peek(input, 1)
		if next == "/" {
			eatCommentLine(input)
			ignored = true
		}
		return token, nil, ignored
	case c == "(":
		return Token{posToken, "LEFT_PAREN", "(", nil}, nil, ignored
	case c == ")":
		return Token{posToken, "RIGHT_PAREN", ")", nil}, nil, ignored
	case c == "{":
		return Token{posToken, "LEFT_BRACE", "{", nil}, nil, ignored
	case c == "}":
		return Token{posToken, "RIGHT_BRACE", "}", nil}, nil, ignored
	case c == "*":
		return Token{posToken, "STAR", "*", nil}, nil, ignored
	case c == ".":
		return Token{posToken, "DOT", ".", nil}, nil, ignored
	case c == ",":
		return Token{posToken, "COMMA", ",", nil}, nil, ignored
	case c == ";":
		return Token{posToken, "SEMICOLON", ";", nil}, nil, ignored
	case c == "+":
		return Token{posToken, "PLUS", "+", nil}, nil, ignored
	case c == "-":
		return Token{posToken, "MINUS", "-", nil}, nil, ignored
	default:
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", line, c)
		return Token{}, errors.New("Error"), ignored
	}
}

func printToken(token Token) {
	fmt.Print(token.Type + " ")
	fmt.Print(token.Token + " ")
	if token.Value != nil {
		fmt.Println(*token.Value)
	} else {
		fmt.Println("null")
	}
}

func printTokenValue(token Token) {
	if token.Type == "EOF" {
		return
	}
	if token.Type == "NUMBER" {
		fmt.Println(*token.Value)
	} else if token.Type == "STRING" {
		fmt.Println(token.Token[1 : len(token.Token)-1])
	} else {
		fmt.Println(token.Token)
	}
}

func advance(input string) (string, bool) {
	posToken++
	if posToken < len(input) {
		return input[posToken : posToken+1], false
	}
	return "", true
}

func peek(input string, pos int) string {
	if posToken+pos < len(input) {
		return input[posToken+pos : posToken+pos+1]
	}
	return ""
}

func eatCommentLine(input string) {
	for {
		c, end := advance(input)
		if c == "\n" || end {
			if c == "\n" {
				posToken--
			}
			break
		}
	}
}

func eatString(input string) (string, error) {
	ret := ""
	open := true
	for {
		c, end := advance(input)
		if c == "\"" || end {
			if c == "\"" {
				open = false
			}
			break
		}
		ret += c
	}
	if open {
		return "", errors.New("unterminated string")
	}
	return ret, nil
}

func eatNumber(input string) (string, error, bool) {
	ret := ""
	decimals := false
	for {
		c, end := advance(input)
		dot := peek(input, 1)
		if end || isWhitespace(c) {
			break
		}
		if dot == "." {
			if decimals {
				ret += c
				return ret, nil, decimals
			}
			next := peek(input, 2)
			if next != "" && isDigit(next) {
				decimals = true
				ret += c
				posToken++
				ret += dot
				posToken++
				ret += next
			} else {
				ret += c
				break
			}
		} else if isDigit(c) {
			ret += c
		} else {
			posToken--
			break
		}
	}
	return ret, nil, decimals
}

func eatIdentifier(input string) string {
	ret := ""
	for {
		c, end := advance(input)
		if end || isWhitespace(c) {
			break
		}
		if isAlpha(c) || isDigit(c) {
			ret += c
		} else {
			posToken--
			break
		}
	}
	return ret
}

func isDigit(c string) bool {
	return c[0] >= '0' && c[0] <= '9'
}

func isAlpha(c string) bool {
	return c[0] >= 'A' && c[0] <= 'Z' || c[0] >= 'a' && c[0] <= 'z' || c[0] == '_'
}

func isWhitespace(c string) bool {
	if c[0] == '\n' {
		line++
	}
	return c[0] == ' ' || c[0] == '\n' || c[0] == '\t'
}
