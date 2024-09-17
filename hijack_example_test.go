package hijack_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/machinebox/graphql"

	"github.com/honestbank/hijack/v2"
	"github.com/honestbank/hijack/v2/handlers"
	"github.com/honestbank/hijack/v2/internal/request"
)

type OperationResult struct {
	Character struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"character"`
}

func ExampleStart() {
	defer hijack.Start(func(r *request.GraphqlRequest) (string, error) {
		return `{
  "data": {
    "character": {
      "id": "11 (Mock)",
      "name": "Albert Einstein (from mock)",
      "status": "Alive (mock)"
    }
  }
}`, nil
	}, "rickandmortyapi.com")()
	client := graphql.NewClient("https://rickandmortyapi.com/graphql")
	request := graphql.NewRequest(`query GetCharacterByID{ character(id:"11"){id, name, status }}`)
	response := OperationResult{}

	_ = client.Run(context.Background(), request, &response)

	fmt.Println(response.Character.ID)
	fmt.Println(response.Character.Name)
	fmt.Println(response.Character.Status)
	//Output:11 (Mock)
	// Albert Einstein (from mock)
	// Alive (mock)
}

func ExampleStart_2() { //nolint:govet
	h := handlers.New()
	h.Set("GetCharacterByID", func(r *request.GraphqlRequest) (string, error) {
		return `{
  "data": {
    "character": {
      "id": "11 (Mock)",
      "name": "Albert Einstein (from mock)",
      "status": "Alive (mock)"
    }
  }
}`, nil
	})
	h.Set("BadOperation", func(r *request.GraphqlRequest) (string, error) {
		return "", errors.New("error")
	})
	defer hijack.Start(h.Handle, "rickandmortyapi.com")()

	client := graphql.NewClient("https://rickandmortyapi.com/graphql")
	request := graphql.NewRequest(`query GetCharacterByID{ character(id:"11"){id, name, status }}`)
	response := OperationResult{}

	_ = client.Run(context.Background(), request, &response)

	fmt.Println(response.Character.ID)
	fmt.Println(response.Character.Name)
	fmt.Println(response.Character.Status)
	//Output:11 (Mock)
	// Albert Einstein (from mock)
	// Alive (mock)
}

func Example_AssertVariablesInRequest() { //nolint:govet
	h := handlers.New()
	h.Set("GetCharacterByID", func(r *request.GraphqlRequest) (string, error) {
		id, err := hijack.GetVariableAs[string](r, "id")
		if err != nil {
			return "", err
		}
		if id == "12" {
			return `
{
  "data": {},
  "errors": [
    {
      "message": "Character with ID 12 not found."
    }
  ]
}`, nil
		}
		return `{
  "data": {
    "character": {
      "id": "11 (Mock)",
      "name": "Albert Einstein (from mock)",
      "status": "Alive (mock)"
    }
  }
}`, nil
	})
	defer hijack.Start(h.Handle, "rickandmortyapi.com")()

	client := graphql.NewClient("https://rickandmortyapi.com/graphql")
	request := graphql.NewRequest(`query GetCharacterByID($id: ID){ character(id:$id){id, name, status }}`)
	request.Var("id", "12")
	response := OperationResult{}

	fmt.Println("Fetching non-existing character...")
	err := client.Run(context.Background(), request, &response)
	fmt.Println(err)

	request = graphql.NewRequest(`query GetCharacterByID($id: ID){ character(id:$id){id, name, status }}`)
	request.Var("id", "11")
	response = OperationResult{}
	fmt.Println("Fetching existing character...")
	err = client.Run(context.Background(), request, &response)
	if err != nil {
		log.Fatal("unexpected error:", err)
	}
	fmt.Println(response.Character.ID)
	fmt.Println(response.Character.Name)
	fmt.Println(response.Character.Status)
	//Output: Fetching non-existing character...
	//graphql: Character with ID 12 not found.
	//Fetching existing character...
	//11 (Mock)
	//Albert Einstein (from mock)
	//Alive (mock)
}
