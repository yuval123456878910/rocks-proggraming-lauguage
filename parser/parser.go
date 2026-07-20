package parser

import (
	"fmt"
	"rocks/lexer"
	"slices"
)

var aviableTypes []string = []string{
	"string", "float", "list", "dict", "bool", "char", "int",
}

func FindNexer(StartPos int, Tokens []lexer.Token, startFind lexer.Token, endFind lexer.Token) (int, error) {
	SkipEnds := 0 // this ment to fix when the code find } but it isnt the end
	for idx, token := range Tokens[StartPos+1:] {
		if token == endFind && SkipEnds <= 0 {
			return StartPos + idx, nil

		} else if token == endFind && SkipEnds > 0 {
			SkipEnds--
		} else if token == startFind {
			SkipEnds++
		}
	}
	return 0, fmt.Errorf("The system coudnt find end for '{'")
}

type NewIdent struct {
	Name    string
	Type    string
	IsConst bool
	Content any
}

type Parser struct {
	Input  []lexer.Token
	Output []any
}

func (p *Parser) Parsing() {
	p.Output = Parse(p.Input)
}

type Function struct {
	Name       string
	Perameters []Parimiter  // type and name
	Returns    []TypeReturn // types
	Body       []any        // all the function code
}

type Parimiter struct {
	Type string
	Name string
}

type TypeReturn struct {
	Type string
}

type Reach struct {
	Path string
}

type Moudle struct {
	Body []any
}

type Return struct {
	Expr any
}

func SearchEnd(Tokens []lexer.Token, pos int) int {
	End := pos
	for End < len(Tokens) && Tokens[End].Type != lexer.NEWLINE {
		End++
	}
	return End
}

func Parse(Tokens []lexer.Token) []any {
	Global_Result := []any{}
	for pos := 0; pos < len(Tokens); pos++ {
		Token := Tokens[pos]
		fmt.Println(Token)
		if pos > 0 && Tokens[pos-1].Type == lexer.IDENTIFIER && Token.Value == "(" && Token.Type == lexer.PUNCTUATOR {

			End, err := FindClose(Tokens, pos+1, "(", ")")
			if err != nil {
				fmt.Println("Error: End of the call funcion wasnt found!")
			}
			Call := CallFunction{Name: Tokens[pos-1].Value}
			InPrunctiotion := false

			var CurrentArg []lexer.Token
			for _, token := range Tokens[pos+1 : End] {
				if token.Type == lexer.PUNCTUATOR && slices.Contains([]string{"{", "}", "(", ")", "{", "}"}, token.Value) {
					InPrunctiotion = true
					pos++
				}

				if token.Type == lexer.PUNCTUATOR && token.Value == "," && !InPrunctiotion {

					TempExpress := Expretion{Tokens: CurrentArg}
					fmt.Println(TempExpress)
					Call.ParimitersInput = append(Call.ParimitersInput, ParseBinding(&TempExpress, 0))
					CurrentArg = []lexer.Token{}
					continue
				}
				CurrentArg = append(CurrentArg, token)

			}
			if len(CurrentArg) > 0 {
				TempExpress := Expretion{Tokens: CurrentArg}
				Call.ParimitersInput = append(Call.ParimitersInput, ParseBinding(&TempExpress, 0))
			}
			Global_Result = append(Global_Result, Call)
			pos = End + 1

		}
		switch Token.Value {
		case "func":
			// syntax: func main(int main) [int] {}
			var Current_Result Function
			Current_Result.Name = Tokens[pos+1].Value // get name of the function like
			if Tokens[pos+2].Value != "(" {
				fmt.Println("Where the decleration '()'")
				return []any{}
			}
			StartToken := lexer.Token{Value: "(", Type: lexer.PUNCTUATOR}
			EndToken := lexer.Token{Value: ")", Type: lexer.PUNCTUATOR}
			index, err := FindNexer(pos+2, Tokens, StartToken, EndToken)
			if err != nil {
				fmt.Println("Error: The perameters didnt end!")
				return nil
			}
			if index+1 > pos {
				for idx, perameters := range Tokens[pos+3 : index+1] {
					if perameters.Type == lexer.KEYWORD && slices.Contains(aviableTypes, perameters.Value) {
						TempPar := Parimiter{Type: perameters.Value, Name: Tokens[pos+4:][idx].Value}
						Current_Result.Perameters = append(Current_Result.Perameters, TempPar)
					}
				}
			}
			if Tokens[index+2].Value == "(" && Tokens[index+2].Type == lexer.PUNCTUATOR {

				StartToken2 := lexer.Token{Value: "(", Type: lexer.PUNCTUATOR}
				EndToken2 := lexer.Token{Value: ")", Type: lexer.PUNCTUATOR}
				index2, err := FindNexer(index+2, Tokens, StartToken2, EndToken2)
				if err != nil {
					fmt.Println("An error had accourd, mayby () didnt closed!")
					return []any{}
				}
				for _, returns := range Tokens[index+3 : index2+1] {
					if returns.Type == lexer.KEYWORD && slices.Contains(aviableTypes, returns.Value) {
						Current_Result.Returns = append(Current_Result.Returns, TypeReturn{Type: returns.Value})
					} else {
						fmt.Println("You included untype as a return parameter")
						return []any{}
					}

				}
			}
			StartBodyIdx := pos
			for !((Tokens[StartBodyIdx].Type == lexer.PUNCTUATOR) && (Tokens[StartBodyIdx].Value == "{")) && StartBodyIdx < len(Tokens)-1 {
				StartBodyIdx++
			}
			StartToken2 := lexer.Token{Value: "{", Type: lexer.PUNCTUATOR}
			EndToken2 := lexer.Token{Value: "}", Type: lexer.PUNCTUATOR}
			EndIdx, err := FindNexer(StartBodyIdx, Tokens, StartToken2, EndToken2)
			if err != nil {
				fmt.Println("Error 105, The function didn't end")
			}
			Current_Result.Body = Parse(Tokens[StartBodyIdx+1 : EndIdx+1])
			Global_Result = append(Global_Result, Current_Result)
			pos = EndIdx + 1
		case "var", "const":

			NewIdentefire := NewIdent{}
			NameToken := Tokens[pos+1]
			if NameToken.Type != lexer.IDENTIFIER {
				fmt.Println("Token isn't an identefire,", NameToken.Value)
				return []any{}
			}
			NewIdentefire.Name = NameToken.Value
			TypeToken := Tokens[pos+2]
			if !slices.Contains(aviableTypes, TypeToken.Value) {
				fmt.Printf("'%s' isn't a type in decleration of value '%s'\n", TypeToken.Value, NameToken.Value)
				return []any{}
			}
			NewIdentefire.Type = TypeToken.Value
			if pos+3 >= len(Tokens) || !(Tokens[pos+3].Value == "=" && Tokens[pos+3].Type == lexer.OPERATOR) {
				Global_Result = append(Global_Result, NewIdentefire)
				continue
			}
			EndOfLine := SearchEnd(Tokens, pos)
			NewExpretion := Expretion{Tokens: Tokens[pos+4 : EndOfLine+1]}
			NewIdentefire.IsConst = Tokens[pos].Value == "const"
			AST := ParseBinding(&NewExpretion, 0)
			NewIdentefire.Content = AST
			pos = EndOfLine

			Global_Result = append(Global_Result, NewIdentefire)
			continue
		case "reach":
			NewReach := Reach{}
			path := Tokens[pos+1]
			if path.Type != lexer.STRING {
				fmt.Println("The reach path isnt a string")
				return []any{}
			}
			NewReach.Path = path.Value
			Global_Result = append(Global_Result, NewReach)
			pos = SearchEnd(Tokens, pos)
		case "return":
			fmt.Println("WOW")
			EndLine := SearchEnd(Tokens, pos)
			NewExpretion := Expretion{Tokens: Tokens[pos+1 : EndLine]}

			NewReturn := Return{Expr: ParseBinding(&NewExpretion, 0)}
			Global_Result = append(Global_Result, NewReturn)
			pos = EndLine

		}

	}
	if len(Global_Result) == 0 {
		return []any{}
	}
	return Global_Result
}

type Expretion struct {
	Tokens []lexer.Token
	Pos    int
}
type Oporation struct {
	Op    string
	Left  any
	Right any
}

type CallFunction struct {
	Name            string
	ParimitersInput []any
}

func (e *Expretion) Next() (lexer.Token, error) {
	e.Pos++
	if e.Pos >= len(e.Tokens) {
		return lexer.Token{}, fmt.Errorf("Limit reached")
	}
	return e.Tokens[e.Pos], nil

}

func IsExpress(token lexer.Token) bool {
	return token.Type == lexer.OPERATOR || token.Type == lexer.IDENTIFIER || token.Type == lexer.LITERAL
}

var Binding map[string]float32 = map[string]float32{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

func RetriveBinding(Value string) float32 {
	v, _ := Binding[Value]
	return v
}

func IsOporator(token lexer.Token) bool {
	return token.Type == lexer.OPERATOR
}

func FindClose(tokens []lexer.Token, startPos int, start string, end string) (int, error) {
	DistanceOfFinish := 0
	for idx, token := range tokens[startPos:] {
		if token.Type == lexer.PUNCTUATOR {
			switch token.Value {
			case start:
				DistanceOfFinish++
			case end:
				if DistanceOfFinish == 0 {
					return startPos + idx, nil
				}
				DistanceOfFinish--

			}
		}
	}
	return 0, fmt.Errorf("No end had bean found!")
}

type List struct {
	Items []any
}

// error in the ParseBinding has to start at 2
func ParseBinding(Express *Expretion, min_bind float32) any {
	Leftside := any(Express.Tokens[Express.Pos])
	fmt.Println(Express.Tokens)
	if Leftside.(lexer.Token).Type == lexer.IDENTIFIER && Express.Pos+1 < len(Express.Tokens) && Express.Tokens[Express.Pos+1].Value == "(" {
		End, err := FindClose(Express.Tokens, Express.Pos+2, "(", ")")
		if err != nil {
			fmt.Println("Error: End of the call funcion wasnt found!")
		}
		fmt.Println()
		Call := CallFunction{Name: Express.Tokens[Express.Pos].Value}
		depth := 0
		var CurrentArg []lexer.Token
		for _, token := range Express.Tokens[Express.Pos+2 : End] {

			if token.Type == lexer.PUNCTUATOR && slices.Contains([]string{"[", "(", "{"}, token.Value) {
				depth++

			} else if token.Type == lexer.PUNCTUATOR && slices.Contains([]string{"]", ")", "}"}, token.Value) {
				depth--

			}

			if token.Type == lexer.PUNCTUATOR && token.Value == "," && depth == 0 {
				TempExpress := Expretion{Tokens: CurrentArg}
				Call.ParimitersInput = append(Call.ParimitersInput, ParseBinding(&TempExpress, 0))
				CurrentArg = []lexer.Token{}
				continue

			}
			CurrentArg = append(CurrentArg, token)

		}
		if len(CurrentArg) > 0 {
			TempExpress := Expretion{Tokens: CurrentArg}
			Call.ParimitersInput = append(Call.ParimitersInput, ParseBinding(&TempExpress, 0))
		}
		Leftside = Call
		Express.Pos = End + 1
	} else if Leftside.(lexer.Token).Value == "(" {
		End, err := FindClose(Express.Tokens, Express.Pos+1, "(", ")")
		if err != nil {
			fmt.Println("No, end to '(' had been found")
			return Leftside
		}
		TempExpress := Expretion{Tokens: Express.Tokens[Express.Pos+1 : End]}
		Leftside = ParseBinding(&TempExpress, 0)
		Express.Pos = End + 1
	} else if Leftside.(lexer.Token).Value == "[" && Leftside.(lexer.Token).Type == lexer.PUNCTUATOR {

		End, err := FindClose(Express.Tokens, Express.Pos+1, "[", "]")
		if err != nil {
			fmt.Println("No, end to '[' had been found")
			return Leftside
		}
		Temped := []lexer.Token{}
		var NewList List
		var InInside int = 0
		for pos := Express.Pos + 1; pos <= End; pos++ {
			if InInside > 0 {
				Temped = append(Temped, Express.Tokens[pos])
				continue
			}
			if Express.Tokens[pos].Value == "[" {
				InInside++
				Temped = append(Temped, Express.Tokens[pos])
				continue
			} else if Express.Tokens[pos].Value == "]" {
				InInside--
				Temped = append(Temped, Express.Tokens[pos])
				continue
			}
			if Express.Tokens[pos].Value == "," && Express.Tokens[pos].Type == lexer.PUNCTUATOR {
				TempExpress := Expretion{Tokens: Temped}
				NewList.Items = append(NewList.Items, ParseBinding(&TempExpress, 0))
				Temped = []lexer.Token{}

			} else {
				Temped = append(Temped, Express.Tokens[pos])
			}
		}
		if len(Temped) > 0 {
			TempExpress := Expretion{Tokens: Temped}
			NewList.Items = append(NewList.Items, ParseBinding(&TempExpress, 0))
		}
		Leftside = NewList
		Express.Pos = End + 1

	} else {
		Express.Next()
	}
	for Express.Pos < len(Express.Tokens) {
		CurrentToken := Express.Tokens[Express.Pos]

		if !IsOporator(CurrentToken) {
			Express.Next()
			continue
		}

		if RetriveBinding(CurrentToken.Value) >= float32(min_bind) {
			Express.Next()
			var Rightside any = ParseBinding(Express, RetriveBinding(CurrentToken.Value)+1)
			Leftside = Oporation{Op: CurrentToken.Value, Right: Rightside, Left: Leftside}
		} else {
			break
		}
	}
	return Leftside
}
