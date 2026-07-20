package interpeter

/* Add:
Add sepport to nested functions: add(add(1,2),1)
*/

import (
	"fmt"
	"math"
	"os"
	"rocks/lexer"
	parser "rocks/parser"
	"slices"
	"strconv"
)

var (
	FuncType        string   = ReturnType(parser.Function{})
	ReachType       string   = ReturnType(parser.Reach{})
	NewIdentType    string   = ReturnType(parser.NewIdent{})
	OporationType   string   = ReturnType(parser.Oporation{})
	LiteralsLisr    []string = []string{lexer.LITERAL, lexer.CHAR, lexer.STRING}
	ApproveSideToOp []string = []string{ReturnType(0), ReturnType(0.0), ReturnType(""), ReturnType(byte(0))}
	IntType         string   = ReturnType(0)
)

type Environment struct {
	FuncMap     map[string]parser.Function
	VariableMap map[string]Ident
	ParseDate   []any
	Output      any
}

func ReturnType(obj any) string {
	return fmt.Sprintf("%T", obj)
}

type Ident struct {
	Value   any
	Name    string
	Type    string
	IsConst bool
}

func CalculateOporator(CalData any, indentMap map[string]Ident) any {
	var Value any
	if ReturnType(CalData) == ReturnType(lexer.Token{}) {
		switch CalData.(lexer.Token).Type {
		case lexer.LITERAL:
			value, err := strconv.ParseFloat(CalData.(lexer.Token).Value, 64)
			if err != nil {
				fmt.Printf("Error converting string to float. Error: %s", err.Error())
				os.Exit(1)
			}
			if math.Mod(value, 1) == 0 {
				return int(value)
			}
			return value
		case lexer.STRING:
			return CalData.(lexer.Token).Value
		case lexer.CHAR:
			return byte(CalData.(lexer.Token).Value[0])
		case lexer.IDENTIFIER:
			TempIdent := CalData.(lexer.Token)
			return indentMap[TempIdent.Value].Value
		}
	}
	op, ok := CalData.(parser.Oporation)
	if !ok {
		fmt.Println("Cant continue because there is not oporation selected!", CalData)
		fmt.Printf("%T\n", CalData)
		os.Exit(1)
	}
	switch op.Op {
	case "+":
		LeftSide := CalculateOporator(op.Left, indentMap)
		RightSide := CalculateOporator(op.Right, indentMap)
		if !slices.Contains(ApproveSideToOp, ReturnType(LeftSide)) && !slices.Contains(ApproveSideToOp, ReturnType(RightSide)) {
			fmt.Println("Cant do None type!", ReturnType(LeftSide), ReturnType(RightSide))
		}

		switch left := LeftSide.(type) {
		case int:
			right, ok := RightSide.(int)
			if ok {
				return left + right
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return int(right2) + left
			}

		case float64:
			right, ok := RightSide.(int)
			if ok {
				return (left) + float64(right)
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return right2 + left
			}
			fmt.Println("Cant add two incompadeple types!")
			os.Exit(1)

		case string:
			right, ok := RightSide.(string)
			if !ok {
				fmt.Println("Cant add to a string, a none string", RightSide)
			}
			return left + right
		}
	case "-":
		LeftSide := CalculateOporator(op.Left, indentMap)
		RightSide := CalculateOporator(op.Right, indentMap)
		if !slices.Contains(ApproveSideToOp, ReturnType(LeftSide)) && !slices.Contains(ApproveSideToOp, ReturnType(RightSide)) {
			fmt.Println("Cant do None type!", ReturnType(LeftSide), ReturnType(RightSide))
		}
		switch left := LeftSide.(type) {
		case int:
			right, ok := RightSide.(int)
			if ok {
				return left - right
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return left - int(right2)
			}
		case float64:
			right, ok := RightSide.(int)
			if ok {
				return (left) - float64(right)
			}
			right2, ok2 := RightSide.(float64)
			if ok2 {
				return left - right2
			}
			fmt.Println("Cant add two incompadeple types!")
			os.Exit(1)

		}
	}
	return Value
}

func (Env *Environment) Interpeter() {
	if Env.VariableMap == nil {
		Env.VariableMap = map[string]Ident{}
	}
	if Env.FuncMap == nil {
		Env.FuncMap = map[string]parser.Function{}
	}
	for _, ParseToken := range Env.ParseDate {
		switch ReturnType(ParseToken) {
		case FuncType:
			FuncTemp := ParseToken.(parser.Function)
			Env.FuncMap[FuncTemp.Name] = FuncTemp

		case NewIdentType:
			IdentTemp := ParseToken.(parser.NewIdent)
			TempEnv := Environment{VariableMap: Env.VariableMap, FuncMap: Env.FuncMap, ParseDate: []any{IdentTemp.Content}}
			TempEnv.Interpeter()
			NewIdent := Ident{Value: TempEnv.Output, Name: IdentTemp.Name, Type: IdentTemp.Type, IsConst: IdentTemp.IsConst}
			Env.VariableMap[NewIdent.Name] = NewIdent
		case OporationType:
			OporationTemp := ParseToken.(parser.Oporation)
			Env.Output = CalculateOporator(OporationTemp, Env.VariableMap)
		case ReturnType(lexer.Token{}):
			token := ParseToken.(lexer.Token)
			if token.Type == lexer.IDENTIFIER {
				val, ok := Env.VariableMap[token.Value]
				if !ok {
					fmt.Println("Couldn't find veruble names:", token.Value)
					os.Exit(1)
				}
				Env.Output = val.Value

			} else {
				Env.Output = CalculateOporator(ParseToken, Env.VariableMap)
			}
		case ReturnType(parser.CallFunction{}):
			TempCall := ParseToken.(parser.CallFunction)

			CallFunc, ok := Env.FuncMap[TempCall.Name]
			if !ok {
				fmt.Println("Error, coudnt find the func:", TempCall.Name)
			}
			CallVarMap := map[string]Ident{}
			for idx, ident := range TempCall.ParimitersInput {
				CallVarMap[CallFunc.Perameters[idx].Name] = Ident{Value: CalculateOporator(ident, Env.VariableMap), Name: CallFunc.Perameters[idx].Name, Type: CallFunc.Perameters[idx].Type, IsConst: false}
			}
			NewEnv := Environment{ParseDate: CallFunc.Body, VariableMap: CallVarMap, FuncMap: Env.FuncMap}
			NewEnv.Interpeter()
			Env.Output = NewEnv.Output
		case ReturnType(parser.Return{}):
			TempReturn := ParseToken.(parser.Return)
			Env.Output = CalculateOporator(TempReturn.Expr, Env.VariableMap)

		}
	}

}
