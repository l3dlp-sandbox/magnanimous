package mg

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type ExpressionContent struct {
	expr     *govaluate.EvaluableExpression
	MarkDown bool
	Text     string
	Location Location
}

type iterableExpression struct {
	array    *govaluate.EvaluableExpression
	path     string
	location Location
}

type fileConsumer func(string) error
type itemConsumer func(interface{}) error

type iterable interface {
	forEach(parameters magParams, fc fileConsumer, ic itemConsumer) error
}

func NewExpression(arg string, location Location, isMarkDown bool, original string) Content {
	expr, err := govaluate.NewEvaluableExpression(arg)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ExpressionContent{expr: expr, MarkDown: isMarkDown, Location: location, Text: original}
}

func NewVariable(arg string, location Location, original string, ctx *WebFileContext) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 2 {
		variable, rawExpr := parts[0], parts[1]
		expr, err := govaluate.NewEvaluableExpression(rawExpr)
		if err != nil {
			log.Printf("WARNING: (%s) Unable to eval (defining %s): %s (%s)",
				location.String(), variable, rawExpr, err.Error())
			return unevaluatedExpression(original)
		}
		v, err := expr.Evaluate(*ctx)
		if err != nil {
			log.Printf("WARNING: (%s) eval failure: %s", location.String(), err.Error())
			return unevaluatedExpression(original)
		}
		(*ctx)[variable] = v
		return nil
	}
	log.Printf("WARNING: (%s) malformed define expression: %s", location.String(), arg)
	return unevaluatedExpression(original)
}

func unevaluatedExpression(original string) Content {
	return &StringContent{Text: fmt.Sprintf("{{%s}}", original)}
}

func asIterable(arg string) (iterable, error) {
	if strings.HasPrefix(arg, "(") && strings.HasSuffix(arg, ")") {
		expr, err := govaluate.NewEvaluableExpression(arg)
		if err != nil {
			return nil, err
		}
		return &iterableExpression{array: expr}, nil
	}
	return &iterableExpression{path: arg}, nil
}

func (e *iterableExpression) forEach(parameters magParams, fc fileConsumer, ic itemConsumer) error {
	if e.array != nil {
		v, err := e.array.Eval(parameters)
		if err != nil {
			return err
		}
		for _, item := range v.([]interface{}) {
			err = ic(item)
			if err != nil {
				return err
			}
		}
	} else {
		dir := ResolveFile(e.path, "source", e.location.Origin)
		f, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, item := range f {
			if !item.IsDir() {
				err = fc(item.Name())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (e *ExpressionContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError {
	r, err := e.expr.Eval(magParams{
		webFiles:       files,
		origin:         e.Location,
		inclusionChain: inclusionChain,
	})
	if err == nil {
		writer.Write([]byte(fmt.Sprintf("%v", r)))
	} else {
		log.Printf("WARNING: (%s) eval failure: %s", e.Location.String(), err.Error())
		writer.Write([]byte(fmt.Sprintf("{{%s}}", e.Text)))
	}
	return nil
}

func (e *ExpressionContent) String() string {
	return fmt.Sprintf("ExpressionContent{%s}", e.Text)
}

func (e *ExpressionContent) IsMarkDown() bool {
	return e.MarkDown
}
