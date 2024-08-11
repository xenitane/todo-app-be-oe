package server

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	xenmw "github.com/xenitane/todo-app-be-oe/internals/middleware"
	"github.com/xenitane/todo-app-be-oe/internals/user"
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
	if claims.Username != c.Param("username") && !claims.IsAdmin {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "you don't have permission",
		}
	}
	user, err := s.db.GetUserByUserName(c.Param("username"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "This user does not exist",
			Internal: err,
		}
	}
	c.JSON(http.StatusOK, user)
	return nil
}

func (s *Server) HandleUpdateUser(c echo.Context) error {
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
	if claims.Username != c.Param("username") && !claims.IsAdmin {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "you don't have permission",
		}
	}
	u, err := s.db.GetUserByUserName(c.Param("username"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  "the user you are trying to modify does not exist",
			Internal: err,
		}
	}
	userUpdateReq := new(user.UserUpdateReq)
	if err := c.Bind(userUpdateReq); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "error binding request body",
			Internal: err,
		}
	}
	flag := false
	if userUpdateReq.FirstName != nil {
		*userUpdateReq.FirstName = strings.TrimSpace(*userUpdateReq.FirstName)
		lenNewFN := len(*userUpdateReq.FirstName)
		if lenNewFN > 3 && lenNewFN < 51 && *userUpdateReq.FirstName != u.FirstName {
			flag = true
			u.FirstName = *userUpdateReq.FirstName
		}
	}
	if userUpdateReq.LastName != nil {
		*userUpdateReq.LastName = strings.TrimSpace(*userUpdateReq.LastName)
		lenNewLN := len(*userUpdateReq.LastName)
		if lenNewLN > 3 && lenNewLN < 51 && *userUpdateReq.LastName != u.LastName {
			flag = true
			u.LastName = *userUpdateReq.LastName
		}
	}
	if userUpdateReq.Password != nil && !u.MatchPassword(*userUpdateReq.Password) {
		lenNewPW := len(*userUpdateReq.Password)
		if lenNewPW < 8 || lenNewPW > 72 {
			return &echo.HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "password length not appropriate",
			}
		}
		if err := u.UpdatePassword(*userUpdateReq.Password); err != nil {
			return &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  "some issue with your password",
				Internal: err,
			}
		}
		flag = true
	}

	if userUpdateReq.IsAdmin != nil {
		if !claims.IsAdmin {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "you are unauthorized for this method",
			}
		}
		if u.Username == claims.Username {
			return &echo.HTTPError{
				Code:    http.StatusTeapot,
				Message: "you cannot demote yourself",
			}
		}
		u.IsAdmin = *userUpdateReq.IsAdmin
		flag = true
	}

	if !flag {
		return &echo.HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "invalid request body",
		}
	}
	if err := s.db.UpadteUser(u); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "internal server error",
			Internal: err,
		}
	}

	c.JSON(http.StatusOK, u)

	return nil
}
