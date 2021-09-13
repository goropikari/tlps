package tlps

// Parser is struct of parser
type Parser struct {
	runtime *Runtime
	tokens  TokenList
	current int
}

// NewParser is constructor of Parser
func NewParser(runtime *Runtime, tokens TokenList) *Parser {
	return &Parser{
		runtime: runtime,
		tokens:  tokens,
		current: 0,
	}
}

// Parse parses given tokens
func (p *Parser) Parse() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(ClassTT) {
		return p.classDeclaration()
	}
	if p.match(FunTT) {
		return p.function("function")
	}
	if p.match(IncludeTT) {
		return p.include()
	}
	if p.match(VarTT) {
		stmt, err := p.varDecralation()
		if err != nil {
			p.synchronize()
			return nil, err
		}
		return stmt, nil
	}
	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
	}
	return stmt, nil
}

func (p *Parser) classDeclaration() (Stmt, error) {
	name, err := p.consume(IdentifierTT, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass *Variable
	if p.match(LeftParenTT) {
		p.consume(IdentifierTT, "Expect superclass name.")
		superclass = NewVariable(p.previous()).(*Variable)
		p.consume(RightParenTT, "Expect ')' after superclass name.")
	}
	_, err = p.consume(ColonTT, "Expect ':' after class name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(NewlineTT, "Expect '\\n' after ':'.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LeftBraceTT, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := make([]*Function, 0)
	for !p.check(RightBraceTT) && !p.isAtEnd() {
		// if p.match(PassTT) {
		// 	p.consume(NewlineTT, "Expect '\\n' after pass")
		// 	continue
		// }

		fun, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, fun.(*Function))
	}

	_, err = p.consume(RightBraceTT, "Expect '}' after class body")

	return NewClass(name, superclass, methods), nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(ForTT) {
		return p.forStatement()
	}
	if p.match(IfTT) {
		return p.ifStatement()
	}
	if p.match(ReturnTT) {
		return p.returnStatement()
	}
	if p.match(WhileTT) {
		return p.whileStatement()
	}
	if p.match(LeftBraceTT) {
		return p.blockStatement(NoneBlock)
	}

	return p.expressionStatement()
}

func (p *Parser) blockStatement(typ BlockType) (Stmt, error) {
	keyword := p.previous() // => '{'
	b, err := p.block()
	if err != nil {
		return nil, err
	}
	return NewBlock(b, keyword, typ), nil
}

func (p *Parser) forStatement() (Stmt, error) {
	var initializer Stmt
	var err error
	if p.match(SemicolonTT) {
		initializer = nil
	} else if p.match(VarTT) {
		initializer, err = p.varDecralation()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if !p.check(SemicolonTT) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SemicolonTT, "Expect ';' after loop condition")
	if err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(ColonTT) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(ColonTT, "Expect ':' after for clauses.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(NewlineTT, "Expect '\\n' after for clauses.")
	if err != nil {
		return nil, err
	}

	keyword, err := p.consume(LeftBraceTT, "Expect '{' for `for` loop body")
	if err != nil {
		return nil, err
	}
	body, err := p.blockStatement(ForBlock)
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = NewBlock([]Stmt{body, NewExpression(increment)}, keyword, ForBlock)
	}

	if condition == nil {
		condition = NewLiteral(true)
	}
	body = NewWhile(condition, body)

	if initializer != nil {
		body = NewBlock([]Stmt{initializer, body}, keyword, ForBlock)
	}

	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(ColonTT, "Expect ':' after if condition.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(NewlineTT, "Expect '\\n' after if condition")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftBraceTT, "Expect '{' for if then block")
	thenBranch, err := p.blockStatement(IfBlock)
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt = nil
	if p.match(ElseifTT) {
		elseBranch, err = p.ifStatement()
		if err != nil {
			return nil, err
		}
		return NewIf(condition, thenBranch, elseBranch), nil
	}
	if p.match(ElseTT) {
		_, err := p.consume(ColonTT, "Expect ':' after else.")
		if err != nil {
			return nil, err
		}
		_, err = p.consume(NewlineTT, "Expect '\\n' after else")
		if err != nil {
			return nil, err
		}
		_, err = p.consume(LeftBraceTT, "Expect '{' for if else block")
		if err != nil {
			return nil, err
		}
		elseBranch, err = p.blockStatement(IfBlock)
		if err != nil {
			return nil, err
		}
	}

	return NewIf(condition, thenBranch, elseBranch), nil
}

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	var value Expr = nil
	var err error
	if !p.check(SemicolonTT) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consumeTerm()
	if err != nil {
		return nil, err
	}
	return NewReturn(keyword, value), nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(ColonTT, "Expect ':' after condition")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(NewlineTT, "Expect '\\n' after condition")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LeftBraceTT, "Expect '{' for while body")
	if err != nil {
		return nil, err
	}
	body, err := p.blockStatement(WhileBlock)
	if err != nil {
		return nil, err
	}

	return NewWhile(condition, body), nil
}

func (p *Parser) varDecralation() (Stmt, error) {
	name, err := p.consume(IdentifierTT, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EqualTT) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consumeTerm()
	if err != nil {
		return nil, err
	}

	return NewVar(name, initializer), nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	// _, err = p.consume(Semicolon, "Expect ';' after expression")
	_, err = p.consumeTerm()
	if err != nil {
		return nil, err
	}

	return NewExpression(expr), nil
}

func (p *Parser) function(kind string) (Stmt, error) {
	if p.match(PassTT) {
		// dummy function
		_, err := p.consumeTerm()
		if err != nil {
			return nil, err
		}

		return NewFunction(
			p.previous(),
			[]*Token{},
			[]Stmt{
				NewExpression(NewLiteral("pass")),
			},
		), nil
	}

	name, err := p.consume(IdentifierTT, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftParenTT, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}
	parameters := make([]*Token, 0)

	if !p.check(RightParenTT) {
		for {
			if len(parameters) >= 255 {
				return nil, p.NewParseError(p.peek(), "Can't have more than 255 parameters.")
			}

			token, err := p.consume(IdentifierTT, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, token)

			if !p.match(CommaTT) {
				break
			}
		}
	}

	_, err = p.consume(RightParenTT, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(ColonTT, "Expect ':' after '('")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(NewlineTT, "Expect '\\n' before "+kind+" body.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LeftBraceTT, "Expected an indented block as "+kind+" body.")
	if err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return NewFunction(name, parameters, body), nil
}

func (p *Parser) include() (Stmt, error) {
	s, err := p.consume(StringTT, "Expect string after include")
	if err != nil {
		return nil, err
	}
	_, err = p.consumeTerm()
	if err != nil {
		return nil, err
	}
	return NewInclude(s), nil
}

func (p *Parser) block() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.check(RightBraceTT) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	p.consume(RightBraceTT, "Expect '}' after block.")
	return statements, nil
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(EqualTT) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if expr.IsType(&Variable{}) {
			name := expr.(*Variable).Name
			return NewAssign(name, value), nil
		} else if expr.IsType(&Get{}) {
			get := expr.(*Get)
			return NewSet(get.Object, get.Name, value), nil
		}

		p.runtime.ErrorTokenMessage(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OrTT) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AndTT) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqualTT, EqualEqualTT) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GreaterTT, GreaterEqualTT, LessTT, LessEqualTT) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MinusTT, PlusTT) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SlashTT, StarTT) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BangTT, MinusTT) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnary(operator, right), nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LeftParenTT) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(DotTT) {
			name, err := p.consume(IdentifierTT, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = NewGet(expr, name)
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := make([]Expr, 0)
	if !p.check(RightParenTT) {
		for {
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}

			if len(arguments) >= 255 {
				return nil, p.NewParseError(p.peek(), "Can't have more than 255 arguments.")
			}

			arguments = append(arguments, expr)

			if !p.match(CommaTT) {
				break
			}
		}
	}

	paren, err := p.consume(RightParenTT, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return NewCall(callee, paren, arguments), nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FalseTT) {
		return NewLiteral(false), nil
	}
	if p.match(TrueTT) {
		return NewLiteral(true), nil
	}
	if p.match(NilTT) {
		return NewLiteral(nil), nil
	}
	if p.match(PassTT) {
		return NewLiteral("pass"), nil
	}
	if p.match(NumberTT, StringTT) {
		return NewLiteral(p.previous().Literal), nil
	}
	if p.match(SuperTT) {
		keyword := p.previous()
		_, err := p.consume(DotTT, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}
		method, err := p.consume(IdentifierTT, "Expect supercrlass method name.")
		if err != nil {
			return nil, err
		}
		return NewSuper(keyword, method), nil
	}
	if p.match(ThisTT) {
		return NewThis(p.previous()), nil
	}
	if p.match(IdentifierTT) {
		return NewVariable(p.previous()), nil
	}
	if p.match(LeftParenTT) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RightParenTT, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return NewGrouping(expr), nil
	}
	if p.match(NewlineTT) {
		return NewLiteral('\n'), nil
	}

	return nil, p.NewParseError(p.peek(), "Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(typ TokenType, message string) (*Token, error) {
	if p.check(typ) {
		return p.advance(), nil
	}

	return nil, p.NewParseError(p.peek(), message)
}

func (p *Parser) consumeTerm() (*Token, error) {
	if p.check(NewlineTT) {
		return p.advance(), nil
	}
	if p.check(SemicolonTT) {
		t := p.advance()
		if p.check(NewlineTT) {
			p.advance()
		}
		return t, nil
	}

	return nil, p.NewParseError(p.peek(), "Expect ';' or '\\n' after expression")
}

func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == typ
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOFTT
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		switch p.previous().Type {
		case SemicolonTT:
			return
		case NewlineTT:
			return
		}
		// if p.previous().Type == Semicolon {
		// 	return
		// }

		switch p.peek().Type {
		case ClassTT:
			return
		case FunTT:
			return
		case VarTT:
			return
		case ForTT:
			return
		case IfTT:
			return
		case WhileTT:
			return
		case ReturnTT:
			return
		}

		p.advance()
	}
}

// NewParseError is constructor of ParseError
func (p *Parser) NewParseError(token *Token, message string) error {
	p.runtime.ErrorTokenMessage(token, message)
	return ParseError.New(token, message)
}
