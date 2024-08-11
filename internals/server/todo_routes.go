package server

import "github.com/labstack/echo/v4"

func (s *Server) RegisterTodoRoutes(g *echo.Group) {
	g.GET("/", nil)
	g.POST("/", nil)
	todoGroup := g.Group("/:todo")
	todoGroup.GET("/", nil)
	todoGroup.PATCH("/", nil)
	todoGroup.DELETE("/", nil)

}
