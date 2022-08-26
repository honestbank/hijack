package graphql

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

	"github.com/honestbank/hijack/internal/generics"
)

func GetOperationName(operation string) (string, error) {
	document, err := parser.Parse(parser.ParseParams{Source: &source.Source{Body: []byte(operation)}})
	if err != nil {
		return "", err
	}
	operationDefinitions := generics.Filter(document.Definitions, func(item ast.Node) bool {
		return item.GetKind() == "OperationDefinition"
	})
	if len(operationDefinitions) != 1 {
		return "", fmt.Errorf("expected exactly 1 OperationDefinition, found %d", len(operationDefinitions))
	}
	name := operationDefinitions[0].(*ast.OperationDefinition).GetName()
	if name == nil {
		return "", errors.New("an operation must have a name")
	}
	return name.Value, nil
}
