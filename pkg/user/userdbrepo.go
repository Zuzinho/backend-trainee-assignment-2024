package user

import (
	"context"
	"database/sql"
)

type UsersDBRepository struct {
	DB *sql.DB
}

func NewUsersDBRepository(db *sql.DB) *UsersDBRepository {
	return &UsersDBRepository{
		DB: db,
	}
}

func (repo *UsersDBRepository) SignUp(ctx context.Context, user *CreateUser) error {
	_, err := repo.DB.ExecContext(ctx, "insert into users (login, password, role) values "+
		"($1, $2, $3)", user.Login, user.Password, user.Role)

	return err
}

func (repo *UsersDBRepository) SignIn(ctx context.Context, login, password string) (Role, error) {
	var role Role

	err := repo.DB.QueryRowContext(ctx, "select role from users where login = $1 and password = $2", login, password).Scan(&role)

	return role, err
}
