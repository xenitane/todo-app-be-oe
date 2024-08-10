package main

import (
	"fmt"

	"github.com/xenitane/todo-app-be-oe/internals/server"
)

func main() {
	server := server.New()

	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("Cannot start the server: %s", err.Error()))
	}
}
