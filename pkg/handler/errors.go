package handler

import "fmt"

type NoRequiredParamError struct {
	key string
}

func newNoRequiredParamError(key string) NoRequiredParamError {
	return NoRequiredParamError{
		key: key,
	}
}

func (err NoRequiredParamError) Error() string {
	return fmt.Sprintf("no required param '%s'", err.key)
}

type NoBannerIDError struct {
}

func (NoBannerIDError) Error() string {
	return "no banner id in path"
}

type IncorrectRoleFromContextError struct {
}

func (IncorrectRoleFromContextError) Error() string {
	return "incorrect role from request context"
}

var (
	NoRequiredParamErr          = NoRequiredParamError{}
	NoBannerIDErr               = NoBannerIDError{}
	IncorrectRoleFromContextErr = IncorrectRoleFromContextError{}
)
