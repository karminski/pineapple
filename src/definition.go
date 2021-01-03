package pineapple 

type Variable struct {
    LineNum int 
    Name    string 
}

type Assignment struct {
    LineNum   int 
    Variable *Variable
    String    string 
}

type Print struct {
    LineNum   int 
    Variable *Variable
}

type Statement interface{}

var _ Statement = (*Print)(nil)
var _ Statement = (*Assignment)(nil)

type SourceCode struct {
    LineNum      int 
    Statements []Statement
}
