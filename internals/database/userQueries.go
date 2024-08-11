package database

import (
	"database/sql"

	"github.com/xenitane/todo-app-be-oe/internals/user"
)

func (s *service) GetAllUsers() ([]*user.User, error) {
	query := `select * from users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	users := []*user.User{}
	for rows.Next() {
		user, err := scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *service) InsertUser(u *user.User) error {
	insertQry := `
		insert into
		users (username,first_name,last_name,password,is_admin)
		values ($1, $2, $3, $4, $5);
		`
	_, err := s.db.Query(
		insertQry,
		u.Username,
		u.FirstName,
		u.LastName,
		u.Password,
		u.IsAdmin,
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetUserByUserName(username string) (*user.User, error) {
	query := `select * from users where username = $1`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanUserRow(rows)
	}
	return nil, sql.ErrNoRows
}

func scanUserRow(rows *sql.Rows) (*user.User, error) {
	user := new(user.User)
	err := rows.Scan(
		&user.UserId,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
	)
	return user, err
}
