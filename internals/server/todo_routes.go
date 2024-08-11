package server

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	xenmw "github.com/xenitane/todo-app-be-oe/internals/middleware"
	"github.com/xenitane/todo-app-be-oe/internals/todo"
)

func (s *Server) RegisterTodoRoutes(g *echo.Group) {
	g.GET("/", s.HandleGetAllTodosOfUser)
	g.POST("/", s.HandleAddTodoForUser)
	todoGroup := g.Group("/:todo")
	todoGroup.GET("/", s.HandleGetTodoByIDForUser)
	todoGroup.PATCH("/", s.HandleUpdateTodoByIDForUser)
	todoGroup.DELETE("/", s.HandleDeleteTodoByIDForUser)
}

func (s *Server) HandleGetAllTodosOfUser(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "jwt missing",
		}
	}
	claims, ok := token.Claims.(*xenmw.JWTCustomClaims)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed to cast claims",
		}
	}
	if claims.Username != c.Param("username") && !claims.IsAdmin {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "you dont have access",
		}
	}
	u, err := s.db.GetUserByUserName(c.Param("username"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "this user does not exist",
			Internal: err,
		}
	}

	todos, err := s.db.GetAllTodosForUser(u.UserId)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "internal server error",
			Internal: err,
		}
	}
	c.JSON(http.StatusOK, todos)
	return nil
}

func (s *Server) HandleAddTodoForUser(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "jwt missing",
		}
	}
	claims, ok := token.Claims.(*xenmw.JWTCustomClaims)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed to cast claims",
		}
	}
	if claims.Username != c.Param("username") {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "You don't have permissions for this method",
		}
	}
	u, err := s.db.GetUserByUserName(c.Param("username"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "this user does not exist",
			Internal: err,
		}
	}
	todoAddReq := new(todo.TodoAddReq)
	if err := c.Bind(todoAddReq); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "error binding request body",
			Internal: err,
		}
	}
	todoAddReq.Title = strings.TrimSpace(todoAddReq.Title)
	todoAddReq.Description = strings.TrimSpace(todoAddReq.Description)
	if err := s.v.Struct(todoAddReq); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "invalid request body",
			Internal: err,
		}
	}
	todo := todo.NewFromAdd(todoAddReq, u.UserId)
	if err := s.db.InsertTodo(todo); err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "internal server error",
			Code:     http.StatusInternalServerError,
		}
	}
	todo.CreatedAt = time.Now()
	c.JSON(http.StatusCreated, todo)
	return nil
}

func (s *Server) HandleGetTodoByIDForUser(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "jwt missing",
		}
	}
	claims, ok := token.Claims.(*xenmw.JWTCustomClaims)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed to cast claims",
		}
	}
	if claims.Username != c.Param("username") && !claims.IsAdmin {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "you dont have access",
		}
	}
	u, err := s.db.GetUserByUserName(c.Param("username"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "this user does not exist",
			Internal: err,
		}
	}
	todoId, err := strconv.ParseInt(c.Param("todo"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Internal: err,
			Message:  "invalid todo id format",
		}
	}
	todo, err := s.db.GetTodoByIDForUser(todoId, u.UserId)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "this user has not todo with this the given id",
			Internal: err,
		}
	}
	c.JSON(http.StatusOK, todo)
	return nil
}

func (s *Server) HandleUpdateTodoByIDForUser(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "jwt missing",
		}
	}
	claims, ok := token.Claims.(*xenmw.JWTCustomClaims)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed to cast claims",
		}
	}
	if claims.Username != c.Param("username") {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "You don't have permissions for this method",
		}
	}
	u, err := s.db.GetUserByUserName(c.Param("username"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "this user does not exist",
			Internal: err,
		}
	}
	todoUpdateReq := new(todo.TodoUpdateReq)
	if err := c.Bind(todoUpdateReq); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "error binding request body",
			Internal: err,
		}
	}
	todoId, err := strconv.ParseInt(c.Param("todo"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Internal: err,
			Message:  "invalid todo id format",
		}
	}
	todo, err := s.db.GetTodoByIDForUser(todoId, u.UserId)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "this user has not todo with this the given id",
			Internal: err,
		}
	}
	flag := false
	if todoUpdateReq.Title != nil {
		*todoUpdateReq.Title = strings.TrimSpace(*todoUpdateReq.Title)
		tlen := len(*todoUpdateReq.Title)
		if tlen > 3 && tlen < 51 && *todoUpdateReq.Title != todo.Title {
			flag = true
			todo.Title = *todoUpdateReq.Title
		}
	}
	if todoUpdateReq.Description != nil {
		*todoUpdateReq.Description = strings.TrimSpace(*todoUpdateReq.Description)
		dlen := len(*todoUpdateReq.Description)
		if dlen < 321 && todo.Description != *todoUpdateReq.Description {
			flag = true
			todo.Description = *todoUpdateReq.Description
		}
	}
	if todoUpdateReq.Status != nil && *todoUpdateReq.Status >= 0 && *todoUpdateReq.Status < 3 && todo.Status != *todoUpdateReq.Status {
		todo.Status = *todoUpdateReq.Status
		flag = true
	}
	if todoUpdateReq.DueDate != nil && !todoUpdateReq.DueDate.Before(time.Now().Round(0)) && reflect.DeepEqual(todo.DueDate, *todoUpdateReq.DueDate) {
		todo.DueDate = *todoUpdateReq.DueDate
		flag = true
	}
	if !flag {
		return &echo.HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "invalid request body",
		}
	}
	if err := s.db.UpdateTodoByIdForUser(todo); err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "internal server error",
			Code:     http.StatusInternalServerError,
		}
	}
	c.JSON(http.StatusCreated, todo)
	return nil
}

func (s *Server) HandleDeleteTodoByIDForUser(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "jwt missing",
		}
	}
	claims, ok := token.Claims.(*xenmw.JWTCustomClaims)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed to cast claims",
		}
	}
	if claims.Username != c.Param("username") && !claims.IsAdmin {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "you dont have access",
		}
	}
	u, err := s.db.GetUserByUserName(c.Param("username"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "this user does not exist",
			Internal: err,
		}
	}
	todoId, err := strconv.ParseInt(c.Param("todo"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Internal: err,
			Message:  "invalid todo id format",
		}
	}
	err = s.db.DeleteTodoByIDForUser(todoId, u.UserId)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "this user has no todo with this the given id",
			Internal: err,
		}
	}

	return nil
}
