package handler

import (
	"avito_hr/pkg/banner"
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type AppHandler struct {
	BannersRepo     banner.BannersRepository
	BannersTempRepo banner.BannersTempRepository
	jsonWorker      *jsonWorker
}

func NewAppHandler(bannersRepo banner.BannersRepository, bannersTempRepo banner.BannersTempRepository) *AppHandler {
	return &AppHandler{
		BannersRepo:     bannersRepo,
		BannersTempRepo: bannersTempRepo,
		jsonWorker:      new(jsonWorker),
	}
}

func (handler *AppHandler) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	tagID, err := handler.intFromQuery(&query, "tag_id", true)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	featureID, err := handler.intFromQuery(&query, "feature_id", true)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	role, err := handler.getRoleFromContext(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	useLastRevision := query.Get("use_last_revision") == "true"

	var content *banner.Content

	if useLastRevision {
		content, err = handler.BannersRepo.GetContent(
			r.Context(),
			tagID,
			featureID,
			role,
		)
		if err != nil {
			if errors.As(err, &sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		content = handler.BannersTempRepo.GetContent(tagID, featureID, role)
	}

	if err = handler.jsonWorker.jsonWrite(w, content); err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) GetBanners(w http.ResponseWriter, r *http.Request) {
	if !handler.checkAdminRole(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	query := r.URL.Query()

	tagID, err := handler.intFromQuery(&query, "tag_id", false)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	featureID, err := handler.intFromQuery(&query, "feature_id", false)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	limit, err := handler.intFromQuery(&query, "tag_id", false)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	offset, err := handler.intFromQuery(&query, "feature_id", false)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	banners, err := handler.BannersRepo.GetBanners(
		r.Context(),
		tagID,
		featureID,
		limit,
		offset,
	)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	if err = handler.jsonWorker.jsonWrite(w, banners); err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	if !handler.checkAdminRole(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	b := new(banner.CreateBanner)

	err := handler.jsonWorker.jsonReadFromRequest(r, b)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	insertedID, err := handler.BannersRepo.CreateBanner(
		r.Context(),
		b,
	)

	resp := postBannerResponse{
		BannerID: insertedID,
	}

	if err = handler.jsonWorker.jsonWrite(w, resp); err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *AppHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	if !handler.checkAdminRole(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)

	bannerID, err := handler.intFromVars(vars, "banner_id")
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	b := new(banner.UpdateBanner)

	err = handler.jsonWorker.jsonReadFromRequest(r, b)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	rowsEff, err := handler.BannersRepo.UpdateBanner(r.Context(), bannerID, b)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	if rowsEff == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AppHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	if !handler.checkAdminRole(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)

	bannerID, err := handler.intFromVars(vars, "banner_id")
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusBadRequest)
		return
	}

	rowsEff, err := handler.BannersRepo.DeleteBanner(r.Context(), bannerID)
	if err != nil {
		handler.jsonWorker.jsonErrorToHTTP(w, err, http.StatusInternalServerError)
		return
	}

	if rowsEff == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
