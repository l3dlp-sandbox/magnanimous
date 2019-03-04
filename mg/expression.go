package mg

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
	"strings"
)

type DefineContent struct {
	Name     string
	Expr     *expression.Expression
	Location *Location
}

type ExpressionContent struct {
	Expr     *expression.Expression
	Text     string
	Location *Location
}

func NewExpression(arg string, location *Location, original string) Content {
	expr, err := expression.ParseExpr(arg)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ExpressionContent{Expr: &expr, Location: location, Text: original}
}

func NewVariable(arg string, location *Location, original string) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 2 {
		variable, rawExpr := parts[0], parts[1]
		expr, err := expression.ParseExpr(rawExpr)
		if err != nil {
			log.Printf("WARNING: (%s) Unable to eval (defining %s): %s (%s)",
				location.String(), variable, rawExpr, err.Error())
			return unevaluatedExpression(original)
		}
		return &DefineContent{Name: variable, Expr: &expr, Location: location}
	}
	log.Printf("WARNING: (%s) malformed define expression: %s", location.String(), arg)
	return unevaluatedExpression(original)
}

func unevaluatedExpression(original string) Content {
	return &StringContent{Text: fmt.Sprintf("{{%s}}", original)}
}

var _ Content = (*ExpressionContent)(nil)

func (e *ExpressionContent) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	r, err := expression.EvalExpr(*e.Expr, magParams{
		stack:    stack,
		webFiles: &files,
	})
	if err == nil {
		_, err = writer.Write([]byte(fmt.Sprintf("%v", r)))
	} else {
		log.Printf("WARNING: (%s) eval failure: %s", e.Location.String(), err.Error())
		_, err = writer.Write([]byte(fmt.Sprintf("{{%s}}", e.Text)))
	}
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func (e *ExpressionContent) String() string {
	return fmt.Sprintf("ExpressionContent{%s}", e.Text)
}

var _ Content = (*DefineContent)(nil)

func (d *DefineContent) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	// DefineContent does not write anything, it just runs an expression and assigns it to a variable
	d.Run(&files, stack)
	return nil
}

func (d *DefineContent) Run(files *WebFilesMap, stack ContextStack) {
	v, err := expression.EvalExpr(*d.Expr, magParams{
		webFiles: files,
		stack:    stack,
	})
	if err != nil {
		log.Printf("WARNING: (%s) define failure: %s", d.Location.String(), err.Error())
	}
	stack.Top().Context.Set(d.Name, v)
}
