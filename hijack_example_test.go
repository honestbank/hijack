package hijack_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/machinebox/graphql"

	"github.com/honestbank/hijack"
	"github.com/honestbank/hijack/handlers"
)

type OperationResult struct {
	Character struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"character"`
}

func ExampleStart() {
	defer hijack.Start(func(operationName string) (string, error) {
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
	h.Set("GetCharacterByID", func(operationName string) (string, error) {
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
	h.Set("BadOperation", func(operationName string) (string, error) {
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
