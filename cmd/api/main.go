package main

import (
	"fmt"
	"os"

	"github.com/xenitane/todo-app-be-oe/internals/server"
)

func main() {

	server := server.New()

	fmt.Printf("starting server at port %v\n", os.Getenv("PORT"))

	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("Cannot start the server: %s", err.Error()))
	}
}
