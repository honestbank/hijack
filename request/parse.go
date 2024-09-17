package request

import (
	"encoding/json"
	"net/http"

	"github.com/honestbank/hijack/v2/internal/graphql"
)

type GraphqlRequest struct {
	OperationName   string                     `json:"operation_name"`
	Query           string                     `json:"query"`
	Variables       map[string]json.RawMessage `json:"variables"`
	OriginalRequest *http.Request
}

func Parse(r *http.Request) (*GraphqlRequest, error) {
	defer r.Body.Close()
	req := GraphqlRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	if req.OperationName == "" {
		req.OperationName, err = graphql.GetOperationName(req.Query)
		if err != nil {
			return nil, err
		}
	}

	return &req, nil
}
