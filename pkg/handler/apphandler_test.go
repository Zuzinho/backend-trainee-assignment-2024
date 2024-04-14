package handler

import (
	"avito_hr/pkg/banner"
	"avito_hr/pkg/user"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserBanner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"get_banner_content"}).
		AddRow([]byte(`{"banner": "content here"}`))
	mock.ExpectQuery("select get_banner_content").
		WithArgs(sql.NullInt64{Int64: 1, Valid: true}, sql.NullInt64{Int64: 1, Valid: true}, user.Role("Admin")).
		WillReturnRows(rows)

	bannersRepo := banner.NewBannersDBRepository(db)
	bannersTempRepo := banner.NewBannersTempMemoryRepository(bannersRepo)
	appHandler := NewAppHandler(bannersRepo, bannersTempRepo)

	req, err := http.NewRequest("GET", "/user_banner?tag_id=1&feature_id=1&use_last_revision=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(context.WithValue(req.Context(), "role", user.RoleAdmin))

	rr := httptest.NewRecorder()

	appHandler.GetUserBanner(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := banner.Content{"banner": "content here"}
	var actual banner.Content
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
