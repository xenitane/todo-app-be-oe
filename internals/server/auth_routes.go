package server

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/xenitane/todo-app-be-oe/internals/middleware"
	"github.com/xenitane/todo-app-be-oe/internals/user"
)

func (s *Server) RegisterAuthRoutes(g *echo.Group) {
	g.POST("/signup/", s.HandleSignup)
	g.POST("/signin/", s.handleSignin)
}

func (s *Server) HandleSignup(c echo.Context) error {
	userReq := new(user.UserSignUpReq)
	if err := c.Bind(userReq); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "error binding request body",
			Internal: err,
		}
	}

	userReq.FirstName = strings.TrimSpace(userReq.FirstName)
	userReq.LastName = strings.TrimSpace(userReq.LastName)
	userReq.Username = strings.TrimSpace(userReq.Username)

	if err := s.v.Struct(userReq); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "invalid request body",
			Internal: err,
		}
	}
	u, err := user.NewFromReg(userReq)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "there is an issue with your password",
			Internal: err,
		}
	}
	if err := s.db.InsertUser(u); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "this username is already taken",
			Internal: err,
		}
	}
	u.CreatedAt = time.Now()
	c.JSON(http.StatusCreated, u)
	return nil
}

func (s *Server) handleSignin(c echo.Context) error {
	userReq := new(user.UserSignInReq)
	if err := c.Bind(userReq); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "error binding request body",
			Internal: err,
		}
	}
	userReq.Username = strings.TrimSpace(userReq.Username)
	if err := s.v.Struct(userReq); nil != err {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "invalid request body",
			Internal: err,
		}
	}
	user, err := s.db.GetUserByUserName(userReq.Username)
	if err != nil {
		return err
	}
	if !user.MatchPassword(userReq.Password) {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "incorrect credentials",
		}
	}

	claims := &middleware.JWTCustomClaims{
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 720)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
	if err != nil {
		return err
	}

	c.Response().Header().Add("x-token-auth", t)
	c.JSON(http.StatusCreated, map[string]any{
		"user":  user,
		"token": t,
	})
	return nil
}
