package main

import (
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
var tokens []Token

func main() {
	posToken = -1

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
		tokens = append(tokens, scanToken(string(fileContents)))
	}
	token := Token{posToken, "EOF", "", nil}
	tokens = append(tokens, token)

	for token := range tokens {
		printToken(tokens[token])
	}
}

func scanToken(input string) Token {
	var token Token
	c := advance(input)
	switch {
	case c == "(":
		token = Token{posToken, "LEFT_PAREN", "(", nil}
	case c == ")":
		token = Token{posToken, "RIGHT_PAREN", ")", nil}
	default:
		token = Token{posToken, "EOF", "", nil}
	}
	return token
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
	posToken = posToken + 1
	return input[posToken : posToken+1]
}
