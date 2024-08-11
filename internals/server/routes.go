package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	xenmw "github.com/xenitane/todo-app-be-oe/internals/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(xenmw.Logger())
	e.Use(xenmw.CORS())

	e.GET("/", s.HiHandler)
	e.GET("/health", s.HealthHandler)

	apiGrp := e.Group("/api")

	s.RegisterAuthRoutes(apiGrp.Group("/auth"))
	s.RegisterUserRoutes(apiGrp.Group("/user", xenmw.JWT()))

	return e
}

func (s *Server) HiHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello there",
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
