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
	RefactType      string   = ReturnType(parser.RefactIdent{})
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

func Evaluate(CalData any, indentMap map[string]Ident, funcMap map[string]parser.Function) any {
	var Value any
	CalType := ReturnType(CalData)
	switch CalType {
	case ReturnType(lexer.Token{}):
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
	case ReturnType(parser.CallFunction{}):
		TempCall := CalData.(parser.CallFunction)

		CallFunc, ok := funcMap[TempCall.Name]
		if !ok {
			fmt.Println("Error, coudnt find the func:", TempCall.Name)
		}
		CallVarMap := map[string]Ident{}
		for idx, ident := range TempCall.ParimitersInput {
			CallVarMap[CallFunc.Perameters[idx].Name] = Ident{Value: Evaluate(ident, indentMap, funcMap), Name: CallFunc.Perameters[idx].Name, Type: CallFunc.Perameters[idx].Type, IsConst: false}
		}
		NewEnv := Environment{ParseDate: CallFunc.Body, VariableMap: CallVarMap, FuncMap: funcMap}
		NewEnv.Interpeter()
		return NewEnv.Output

	}
	op, ok := CalData.(parser.Oporation)
	if !ok {
		fmt.Println("Cant continue because there is not oporation selected!", CalData)
		fmt.Printf("%T\n", CalData)
		os.Exit(1)
	}
	switch op.Op {
	case "+":
		LeftSide := Evaluate(op.Left, indentMap, funcMap)
		RightSide := Evaluate(op.Right, indentMap, funcMap)
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
		LeftSide := Evaluate(op.Left, indentMap, funcMap)
		RightSide := Evaluate(op.Right, indentMap, funcMap)
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
			Env.Output = Evaluate(OporationTemp, Env.VariableMap, Env.FuncMap)

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
				Env.Output = Evaluate(ParseToken, Env.VariableMap, Env.FuncMap)
			}
		case ReturnType(parser.CallFunction{}):
			TempCall := ParseToken.(parser.CallFunction)

			CallFunc, ok := Env.FuncMap[TempCall.Name]
			if !ok {
				fmt.Println("Error, coudnt find the func:", TempCall.Name)
				os.Exit(1)
			}
			CallVarMap := map[string]Ident{}
			for idx, ident := range TempCall.ParimitersInput {
				CallVarMap[CallFunc.Perameters[idx].Name] = Ident{Value: Evaluate(ident, Env.VariableMap, Env.FuncMap), Name: CallFunc.Perameters[idx].Name, Type: CallFunc.Perameters[idx].Type, IsConst: false}
			}
			NewEnv := Environment{ParseDate: CallFunc.Body, VariableMap: CallVarMap, FuncMap: Env.FuncMap}
			NewEnv.Interpeter()
			Env.Output = NewEnv.Output

		case ReturnType(parser.Return{}):
			TempReturn := ParseToken.(parser.Return)
			Env.Output = Evaluate(TempReturn.Expr, Env.VariableMap, Env.FuncMap)
		case RefactType:
			// dont mess: Env.VariableMap[NewIdent.Name] = NewIdent
			TempRefact := ParseToken.(parser.RefactIdent)
			IdentGet, ok := Env.VariableMap[TempRefact.Name]
			if !ok {
				fmt.Printf("Coudnt find a veruble named: %s\n", TempRefact.Name)
				os.Exit(1)
			}
			if IdentGet.IsConst {
				fmt.Println("Cant edit a const veruble: ", TempRefact.Name)
				os.Exit(1)
			}
			NewEnv := Environment{VariableMap: Env.VariableMap, FuncMap: Env.FuncMap, ParseDate: []any{TempRefact.Content}}
			NewEnv.Interpeter()
			var tempIdent Ident = Env.VariableMap[IdentGet.Name]
			tempIdent.Value = NewEnv.Output
			Env.VariableMap[IdentGet.Name] = tempIdent
		}
	}

}
