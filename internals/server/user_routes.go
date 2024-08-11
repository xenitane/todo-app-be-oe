package server

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	xenmw "github.com/xenitane/todo-app-be-oe/internals/middleware"
)

func (s *Server) RegisterUserRoutes(g *echo.Group) {
	g.GET("/", s.HandleAllUsers)
	userGroup := g.Group("/:username")
	userGroup.GET("/", s.HandleUserByUserName)
	userGroup.PATCH("/", s.HandleUpdateUser)
	s.RegisterTodoRoutes(userGroup.Group("/todo"))
}

func (s *Server) HandleAllUsers(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "JWT MISSING",
		}
	}
	claims, ok := token.Claims.(*xenmw.JWTCustomClaims)
	if !ok {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "failed to cast claims",
		}
	}
	if !claims.IsAdmin {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "you are not admin",
		}
	}
	users, err := s.db.GetAllUsers()
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, users)
	return nil
}

func (s *Server) HandleUserByUserName(c echo.Context) error {
	return nil
}

func (s *Server) HandleUpdateUser(c echo.Context) error {
	return nil
}
