package request

import (
	"encoding/json"
	"net/http"

	"github.com/honestbank/hijack/internal/graphql"
)

type Request struct {
	Query string `json:"query"`
}

func GetOperation(r *http.Request) (string, error) {
	defer r.Body.Close()
	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return "", err
	}

	return graphql.GetOperationName(req.Query)
}
