package user

import "fmt"

type UnknownRoleError struct {
	role Role
}

func newUnknownRoleError(role Role) UnknownRoleError {
	return UnknownRoleError{
		role: role,
	}
}

func (err UnknownRoleError) Error() string {
	return fmt.Sprintf("unknown role '%s'", err.role)
}

var (
	UnknownRoleErr = UnknownRoleError{}
)
