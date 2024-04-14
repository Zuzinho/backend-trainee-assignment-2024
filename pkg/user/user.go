package user

import "context"

type Role string

const (
	RoleUser  Role = "User"
	RoleAdmin Role = "Admin"
)

func (role Role) Valid() error {
	switch role {
	case RoleUser, RoleAdmin:
		return nil
	default:
		return newUnknownRoleError(role)
	}
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     Role   `json:"role" validate:"oneof=User Admin"`
}

type CreateUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     *Role  `json:"role,omitempty" validate:"oneof=User Admin"`
}

type LoginUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UsersRepository interface {
	SignUp(ctx context.Context, user *CreateUser) error
	SignIn(ctx context.Context, login, password string) (Role, error)
}
