package handler

import (
	"avito_hr/pkg/user"
	"database/sql"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func (handler *AppHandler) intFromQuery(query *url.Values, key string, required bool) (sql.NullInt64, error) {
	value := query.Get(key)

	if value == "" {
		if required {
			return sql.NullInt64{}, newNoRequiredParamError(key)
		}

		return sql.NullInt64{}, nil
	}

	dig, err := strconv.Atoi(value)

	return sql.NullInt64{Int64: int64(dig), Valid: true}, err
}

func (handler *AppHandler) intFromVars(vars map[string]string, key string) (int64, error) {
	value, ok := vars[key]
	if !ok {
		return 0, NoBannerIDErr
	}

	dig, err := strconv.Atoi(value)

	return int64(dig), err
}

func (handler *AppHandler) checkAdminRole(r *http.Request) bool {
	val := r.Context().Value("role")

	role, ok := val.(user.Role)

	return ok && role == user.RoleAdmin
}

func (handler *AppHandler) getRoleFromContext(r *http.Request) (user.Role, error) {
	val := r.Context().Value("role")

	role, ok := val.(user.Role)
	if !ok {
		return "", IncorrectRoleFromContextErr
	}

	return role, nil
}

type jsonWorker struct {
}

func (worker *jsonWorker) jsonReadFromRequest(r *http.Request, v any) error {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	log.WithFields(log.Fields{
		"Body": string(body),
	}).Info("Trying read request body...")
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func (worker *jsonWorker) jsonErrorToHTTP(w http.ResponseWriter, err error, status int) {
	log.WithFields(log.Fields{
		"Error": err.Error(),
	}).Error("Catching error...")

	w.Header().Set("Content-Type", "application/json")
	http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err.Error()), status)
}

func (worker *jsonWorker) jsonWrite(w http.ResponseWriter, value any) error {
	log.WithFields(log.Fields{
		"Written value": value,
	}).Info("Writing value to response writer...")

	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(buf)

	return err
}
