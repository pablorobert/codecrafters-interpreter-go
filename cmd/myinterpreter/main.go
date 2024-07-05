package main

import (
	"errors"
	"fmt"
	"os"
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

func main() {
	posToken = -1
	line = 1
	hasError = false

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
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
		newToken, err := scanToken(string(fileContents))
		if err != nil {
			hasError = true
		} else {
			tokens = append(tokens, newToken)
		}
	}
	token := Token{posToken, "EOF", "", nil}
	tokens = append(tokens, token)

	for token := range tokens {
		printToken(tokens[token])
	}

	if hasError {
		os.Exit(65)
	}
}

func scanToken(input string) (Token, error) {
	//var token Token
	c := advance(input)
	switch {
	case c == "(":
		return Token{posToken, "LEFT_PAREN", "(", nil}, nil
	case c == "=":
		token := Token{posToken, "EQUAL", "=", nil}
		next := peek(input)
		if next == "=" {
			token = Token{posToken, "EQUAL_EQUAL", "==", nil}
			posToken++
		}
		return token, nil
	case c == ")":
		return Token{posToken, "RIGHT_PAREN", ")", nil}, nil
	case c == "{":
		return Token{posToken, "LEFT_BRACE", "{", nil}, nil
	case c == "}":
		return Token{posToken, "RIGHT_BRACE", "}", nil}, nil
	case c == "*":
		return Token{posToken, "STAR", "*", nil}, nil
	case c == ".":
		return Token{posToken, "DOT", ".", nil}, nil
	case c == ",":
		return Token{posToken, "COMMA", ",", nil}, nil
	case c == ";":
		return Token{posToken, "SEMICOLON", ";", nil}, nil
	case c == "+":
		return Token{posToken, "PLUS", "+", nil}, nil
	case c == "-":
		return Token{posToken, "MINUS", "-", nil}, nil
	default:
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", line, c)
		return Token{}, errors.New("Error")
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

func advance(input string) string {
	posToken++
	return input[posToken : posToken+1]
}

func peek(input string) string {
	if posToken+1 < len(input) {
		return input[posToken+1 : posToken+2]
	}
	return ""
}
