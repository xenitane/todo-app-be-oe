package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/xenitane/todo-app-be-oe/internals/database"
)

type Server struct {
	port int
	v    *validator.Validate
	db   database.Service
}

func New() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		v:    validator.New(),
		db:   database.New(),
	}

	NewServer.v.RegisterValidation("not-stale", validateDateNotStale)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func validateDateNotStale(fl validator.FieldLevel) bool {
	ts, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return !ts.Before(time.Now().Round(0))
}
