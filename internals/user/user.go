package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserSignUpReq struct {
	Username  string `json:"username" validate:"required,min=5,max=20,alphanum"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
	FirstName string `json:"firstName" validate:"required,min=4,max=50,alpha"`
	LastName  string `json:"lastName" validate:"required,min=4,max=50,alpha"`
}

type UserSignInReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateReq struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Password  *string `json:"password"`
	IsAdmin   *bool   `json:"isAdmin"`
}

type User struct {
	UserId    int64     `json:"-"`
	Password  string    `json:"-"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewFromReg(u *UserSignUpReq) (*User, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 7)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Password:  string(hashedPasswordBytes),
	}, nil
}

func (u *User) MatchPassword(password string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) UpdatePassword(password string) error {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return err
	}
	u.Password = string(hashedPasswordBytes)
	return nil
}
