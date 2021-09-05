package tlps

import (
	"bytes"
	"fmt"
	"os"
)

// Runtime is struct of Runtime
type Runtime struct {
	HadError        bool
	HadRuntimeError bool
}

// NewRuntime is constructor of Runtime
func NewRuntime() *Runtime {
	return &Runtime{HadError: false, HadRuntimeError: false}
}

// Run runs script
func (r *Runtime) Run(source *bytes.Buffer) {

	scanner := NewScanner(r, source)
	tokens := scanner.ScanTokens()
	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }

	parser := NewParser(r, tokens)
	expression := parser.Parse()

	// Stop if there was a syntax error
	if r.HadError {
		return
	}

	// fmt.Println(NewAstPrinter().Print(expression))
	interpreter := NewInterpreter(r)
	val, _ := interpreter.Interpret(expression)
	fmt.Println(val)
}

// ErrorMessage prints error massage at stderr
func (r *Runtime) ErrorMessage(line int, message string) {
	r.report(line, "", message)
}

// ErrorTokenMessage prints error message at stderr
func (r *Runtime) ErrorTokenMessage(token *Token, message string) {
	if token.Type == EOF {
		r.report(token.Line, " at end", message)
	} else {
		r.report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

// Report prints error masseg at stderr
func (r *Runtime) report(line int, where string, message string) {
	fmt.Fprintln(os.Stderr, "[line "+fmt.Sprint(line)+"] Error"+where+": "+message)
	r.HadError = true
}

func (r *Runtime) RuntimeError(err error) {
	e := err.(*CustomError)
	fmt.Fprint(os.Stderr, err.Error()+"\n[line "+fmt.Sprint(e.Token.Line)+"]")
	r.HadRuntimeError = true
}
