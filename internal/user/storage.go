package user

import (
	"database/sql"
	"errors"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) createUser(userName, email string, isAdmin bool) (User, error) {
	user := User{
		Username: userName,
		Email:    email,
		IsAdmin:  isAdmin,
	}
	statement := `insert into users(username, email, isAdmin) values($1, $2, $3);`
	_, err := s.db.Exec(statement, user.Username, user.Email, user.IsAdmin)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *UserStorage) IsUserWithEmailPresent(email string) (bool, error) {
	_, err := s.getUserByEmail(email)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *UserStorage) getUserByEmail(email string) (User, error) {
	var user User
	statement := `select * from users where email = $1;`
	row := s.db.QueryRow(statement, email)
	if err := row.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt,
		&user.Username, &user.Email, &user.IsAdmin); err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("USER_NOT_FOUND")
		}
		return User{}, err
	}
	return user, nil
}
