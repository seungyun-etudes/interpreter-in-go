package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

type Type string

const (
	INTEGER_OBJECT      = "INTEGER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NULL_OBJECT         = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return INTEGER_OBJECT
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN_OBJECT
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() Type {
	return NULL_OBJECT
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() Type {
	return RETURN_VALUE_OBJECT
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ERROR_OBJECT
}

func (e *Error) Inspect() string {
	return "ERROR :" + e.Message
}

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (f *Function) Type() Type {
	return FUNCTION_OBJECT
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	var params []string

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("function")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
