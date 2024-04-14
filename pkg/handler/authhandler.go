package handler

import (
	"avito_hr/pkg/session"
	"avito_hr/pkg/user"
	"database/sql"
	"errors"
	"github.com/go-playground/validator"
	"net/http"
)

type AuthHandler struct {
	UsersRepo       user.UsersRepository
	SessionsManager session.SessionsPacker
	jsonWorker      *jsonWorker
	validator       *validator.Validate
}

func NewAuthHandler(usersRepo user.UsersRepository, sessionManager session.SessionsPacker) *AuthHandler {
	return &AuthHandler{
		UsersRepo:       usersRepo,
		SessionsManager: sessionManager,
		jsonWorker:      new(jsonWorker),
		validator:       validator.New(),
	}
}

func (handler *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	u := new(user.CreateUser)

	err := handler.jsonWorker.jsonReadFromRequest(r, u)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	if u.Role == nil {
		role := user.RoleUser
		u.Role = &role
	}

	err = handler.validator.Struct(u)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	err = handler.UsersRepo.SignUp(r.Context(), u)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	token, err := handler.SessionsManager.Pack(*u.Role)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	resp := postUserResponse{
		Token: *token,
	}

	if err = handler.jsonWorker.jsonWrite(w, resp); err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	u := new(user.LoginUser)

	err := handler.jsonWorker.jsonReadFromRequest(r, u)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	role, err := handler.UsersRepo.SignIn(r.Context(), u.Login, u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
			return
		}

		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	token, err := handler.SessionsManager.Pack(role)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	resp := postUserResponse{
		Token: *token,
	}

	if err = handler.jsonWorker.jsonWrite(w, resp); err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
