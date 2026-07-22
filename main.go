package main

import (
	"fmt"
	"os"
	interpeter "rocks/interpeter"
	lexer "rocks/lexer"
	parser "rocks/parser"
)

func main() {
	fmt.Println("START: ")
	var L1 lexer.Lexer
	Data, _ := os.ReadFile("testFile.ro")
	L1.Input = string(Data)
	fmt.Println(L1.Input)

	L1.CorrectBackSlash()
	L1.LexerAll()

	parser := parser.Parser{}
	parser.Input = L1.Tokens
	parser.Parsing()
	fmt.Println(parser.Output...)
	NewEnv := interpeter.Environment{}
	NewEnv.ParseDate = parser.Output
	NewEnv.Interpeter()
	fmt.Println(NewEnv.VariableMap)

}
