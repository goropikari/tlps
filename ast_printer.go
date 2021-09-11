package tlps

import (
	"bytes"
	"fmt"
	"strings"
)

// AstPrinter is struct of ast printer
type AstPrinter struct{}

// NewAstPrinter is constructor of AstPrinter
func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

// Print prints given statements ast
func (ap *AstPrinter) Print(stmts []Stmt) (string, error) {
	val, err := stmts[0].Accept(ap)
	return val.(string), err
}

func (ap *AstPrinter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	return ap.parenthesizeExpr(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap *AstPrinter) visitCallExpr(expr *Call) (interface{}, error) {
	callee, _ := ap.parenthesizeExpr("callee", expr.Callee)
	args := make([]string, 0)
	for _, v := range expr.Arguments {
		arg, _ := ap.parenthesizeExpr("arg", v)
		args = append(args, arg)
	}

	return callee + strings.Join(args, " "), nil
}

func (ap *AstPrinter) visitGroupingExpr(expr *Grouping) (interface{}, error) {
	return ap.parenthesizeExpr("group", expr.Expression)
}

func (ap *AstPrinter) visitLiteralExpr(expr *Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (ap *AstPrinter) visitLogicalExpr(expr *Logical) (interface{}, error) {
	return ap.parenthesizeExpr(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap *AstPrinter) visitUnaryExpr(expr *Unary) (interface{}, error) {
	return ap.parenthesizeExpr(expr.Operator.Lexeme, expr.Right)
}

func (ap *AstPrinter) visitAssignExpr(expr *Assign) (interface{}, error) {
	return ap.parenthesizeExpr("assign "+expr.Name.Lexeme, expr.Value)
}

func (ap *AstPrinter) visitVariableExpr(expr *Variable) (interface{}, error) {
	return ap.parenthesizeExpr("variable", NewLiteral(expr.Name.Lexeme))
}

func (ap *AstPrinter) visitBlockStmt(b *Block) (interface{}, error) {
	body := make([]string, 0)
	for _, stmt := range b.Statements {
		s, err := ap.parenthesizeStmt("block body", stmt)
		if err != nil {
			return "", err
		}
		body = append(body, s)
	}
	return "(block " + strings.Join(body, " ") + ")", nil
}

func (ap *AstPrinter) visitExpressionStmt(e *Expression) (interface{}, error) {
	return e.Expression.Accept(ap)
}

func (ap *AstPrinter) visitFunctionStmt(f *Function) (interface{}, error) {
	params := make([]string, 0)
	for _, v := range f.Params {
		params = append(params, v.Lexeme)
	}
	stmts := make([]string, 0)
	for _, stmt := range f.Body {
		s, err := stmt.Accept(ap)
		if err != nil {
			return "nil", err
		}
		stmts = append(stmts, "("+s.(string)+")")
	}

	return "(function (args (" + strings.Join(params, ", ") + ")) (body " + strings.Join(stmts, " ") + "))", nil
}

func (ap *AstPrinter) visitIfStmt(i *If) (interface{}, error) {
	cond, err := ap.parenthesizeExpr("cond", i.Condition)
	if err != nil {
		return "", nil
	}
	thenBranch, err := ap.parenthesizeStmt("thenBranch", i.ThenBranch)
	if err != nil {
		return "", nil
	}
	var elseBranch string
	if i.ElseBranch != nil {
		elseBranch, err = ap.parenthesizeStmt("elseBranch", i.ElseBranch)
		if err != nil {
			return "", err
		}
	}

	return "(if " + cond + " " + thenBranch + " " + elseBranch + ")", nil
}

func (ap *AstPrinter) visitPrintStmt(p *Print) (interface{}, error) {
	return ap.parenthesizeExpr("print", p.Expression)
}

func (ap *AstPrinter) visitReturnStmt(r *Return) (interface{}, error) {
	expr, err := ap.parenthesizeExpr(r.Keyword.Lexeme, r.Value)
	if err != nil {
		return "nil", err
	}

	return expr, nil
}

func (ap *AstPrinter) visitWhileStmt(p *While) (interface{}, error) {
	cond, err := ap.parenthesizeExpr("cond", p.Condition)
	if err != nil {
		return "", err
	}
	body, err := ap.parenthesizeStmt("body", p.Body)
	if err != nil {
		return "", nil
	}
	return "(while " + cond + " " + body + ")", nil
}

func (ap *AstPrinter) visitVarStmt(v *Var) (interface{}, error) {
	initializer, err := ap.parenthesizeExpr("init", v.Initializer)
	if err != nil {
		return "", nil
	}
	return "(declare " + v.Name.Lexeme + " " + initializer + ")", nil
}

func (ap *AstPrinter) parenthesizeExpr(name string, exprs ...Expr) (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString("(" + name)
	for _, expr := range exprs {
		buf.WriteString(" ")
		s, _ := expr.Accept(ap)
		buf.WriteString(s.(string))
	}
	buf.WriteString(")")

	return buf.String(), nil
}

func (ap *AstPrinter) parenthesizeStmt(name string, stmts ...Stmt) (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString("(" + name)
	for _, stmt := range stmts {
		buf.WriteString(" ")
		s, _ := stmt.Accept(ap)
		buf.WriteString(s.(string))
	}
	buf.WriteString(")")

	return buf.String(), nil
}
