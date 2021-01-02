package pineapple

import (
    "errors"
)


// Name ::= [_A-Za-z][_0-9A-Za-z]*
func parseName(lexer *Lexer) (string, error) {
	_, name := lexer.NextTokenIs(TOKEN_NAME)
	return name, nil
}

// String ::= '"' '"' Ignored | '"' StringCharacter '"' Ignored
func parseString(lexer *Lexer) (string, error) {
	str := "" 
	switch lexer.LookAhead() {
	case TOKEN_DUOQUOTE:
		return str, nil 
	case TOKEN_QUOTE:
		lexer.NextTokenIs(TOKEN_QUOTE)
		str = lexer.scanBeforeToken(tokenNameMap[TOKEN_QUOTE])
		lexer.NextTokenIs(TOKEN_QUOTE)
		return str, nil
	default:
		return "", errors.New("parseString(): not a string.")
	}
} 

// Variable ::= "$" Name Ignored
func parseVariable(lexer *Lexer) (*Variable, error) {
	var variable Variable
	var err      error 
	
	variable.LineNum = lexer.GetLineNum()
	lexer.NextTokenIs(TOKEN_VAR_PREFIX)
	if variable.Name, err = parseName(lexer); err != nil {
		return nil, err
	}
	return &variable, nil
}


// Assignment  ::= Variable Ignored "=" Ignored String Ignored
func parseAssignment(lexer *Lexer) (*Assignment, error) {
	var assignment Assignment
	var err 	   error
	
	assignment.LineNum = lexer.GetLineNum()
	if assignment.Variable, err = parseVariable(lexer); err != nil {
		return nil, err
	}
	lexer.NextTokenIs(TOKEN_EQUAL)
	if assignment.String, err = parseString(lexer); err != nil {
		return nil, err
	}
	return &assignment, nil
}

// Print ::= "print" "(" Ignored Variable Ignored ")" Ignored
func parsePrint(lexer *Lexer) (*Print, error) {
	var print Print 
	var err   error

	print.LineNum = lexer.GetLineNum()
	lexer.NextTokenIs(TOKEN_PRINT)
	lexer.NextTokenIs(TOKEN_LEFT_PAREN)
	if print.Variable, err = parseVariable(lexer); err != nil {
		return nil, err
	}
	lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
	return &print, nil
}


// Statement ::= Print | Assignment
func parseStatements(lexer *Lexer) ([]Statement, error) {
	var statements []Statement 
	
	for !isSourceCodeEnd(lexer.LookAhead()) {
		var statement Statement
		var err       error
		if statement, err = parseStatement(lexer); err != nil {
			return nil, err 
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

func parseStatement(lexer *Lexer) (Statement, error) {
	switch lexer.LookAhead() {
	case TOKEN_PRINT:
		return parsePrint(lexer)
	case TOKEN_VAR_PREFIX:
		return parseAssignment(lexer)
	default:
		return nil, errors.New("parseStatement(): unknown Statement.")
	}
}

// SourceCode ::= Statement+ 
func parseSourceCode(lexer *Lexer) (*SourceCode, error) {
	var sourceCode SourceCode
	var err 	   error

	sourceCode.LineNum = lexer.GetLineNum()
	if sourceCode.Statements, err = parseStatements(lexer); err != nil {
		return nil, err
	}
	return &sourceCode, nil
}

func isSourceCodeEnd(token int) bool {
    if token == TOKEN_EOF {
        return true
    }
    return false
}

func parse(code string) (*SourceCode, error) {
	var sourceCode *SourceCode
	var err 	    error 

	lexer := NewLexer(code)
	if sourceCode, err = parseSourceCode(lexer); err != nil {
		return nil, err 
	}
	lexer.NextTokenIs(TOKEN_EOF)
	return sourceCode, nil
}