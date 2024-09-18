package handlers_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/honestbank/hijack/v2/handlers"
	"github.com/honestbank/hijack/v2/request"
)

func ExampleNew() {
	cb := func(r *request.GraphqlRequest) (string, error) {
		return fmt.Sprintf("Hello from operation %s.\nQuery: %s.\nVariables: %s", r.OperationName, r.Query, r.Variables), nil
	}
	h := handlers.New().
		Set("Operation1", cb).
		Set("Operation2", cb).
		Set("HeroName", cb)

	rawRequestBody := `{
  "query" : "query HeroName($episode: Episode) { hero(episode: $episode) { name } }",
  "variables" : {
   "episode": "JEDI"
  },
  "operationName" : "HeroName"
}`
	reqRaw, err := http.NewRequest("", "", bytes.NewBufferString(rawRequestBody))
	if err != nil {
		log.Fatal(err)
	}
	req, err := request.Parse(reqRaw)
	if err != nil {
		log.Fatal(err)
	}
	res, err := h.Handle(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	// Output: Hello from operation HeroName.
	//Query: query HeroName($episode: Episode) { hero(episode: $episode) { name } }.
	//Variables: map[episode:"JEDI"]
}
