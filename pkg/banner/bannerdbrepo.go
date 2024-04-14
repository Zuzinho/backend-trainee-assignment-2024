package banner

import (
	"avito_hr/pkg/user"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"strings"
)

type BannersDBRepository struct {
	DB *sql.DB
}

func NewBannersDBRepository(db *sql.DB) *BannersDBRepository {
	return &BannersDBRepository{
		DB: db,
	}
}

func (repo *BannersDBRepository) GetContent(ctx context.Context, tagID, featureID sql.NullInt64, role user.Role) (*Content, error) {
	content := new(Content)

	err := repo.DB.QueryRowContext(ctx, "select get_banner_content from get_banner_content($1, $2, $3)",
		tagID, featureID, role).Scan(content)

	return content, err
}

func (repo *BannersDBRepository) GetBanners(ctx context.Context, tagID, featureID, limit, offset sql.NullInt64) (*Banners, error) {
	rows, err := repo.DB.QueryContext(ctx,
		fmt.Sprintf("select banner_id, feature_id, tag_ids, "+
			"content, is_active, created_at, updated_at from get_banners($1, $2, $3, $4)"),
		limit, offset, tagID, featureID,
	)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	banners := make(Banners, 0)
	for rows.Next() {
		banner := &Banner{}

		var tagIDS pq.Int64Array

		err = rows.Scan(&banner.BannerID, &banner.FeatureID, &tagIDS, &banner.Content,
			&banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
		if err != nil {
			log.Error(err)
			continue
		}

		banner.TagIDs = tagIDS

		banners.Append(banner)
	}

	return &banners, nil
}

func (repo *BannersDBRepository) CreateBanner(ctx context.Context, banner *CreateBanner) (int64, error) {
	var insertedID int64

	err := repo.DB.QueryRowContext(ctx,
		"select insert_banner from insert_banner($1, $2, $3, $4)",
		banner.FeatureID, pq.Array(banner.TagIDs),
		banner.Content, banner.IsActive,
	).Scan(&insertedID)

	return insertedID, err
}

func (repo *BannersDBRepository) UpdateBanner(ctx context.Context, bannerID int64, banner *UpdateBanner) (int64, error) {
	argNumber := 2
	args := []any{bannerID}
	sets := make([]string, 0)
	var contentUpdateQuery, tagsUpdateQuery, updateQuery string

	if banner.FeatureID != nil {
		sets = append(sets, fmt.Sprintf("feature_id = $%d", argNumber))
		args = append(args, banner.FeatureID)
		argNumber++
	}

	if banner.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", argNumber))
		args = append(args, banner.IsActive)
		argNumber++
	}

	if banner.Content != nil {
		contentUpdateQuery = fmt.Sprintf("call update_banner_content($1, $%d);", argNumber)
		args = append(args, banner.Content)
		argNumber++
	}

	if banner.TagIDs != nil {
		tagsUpdateQuery = fmt.Sprintf("call update_banner_tags($1, $%d);", argNumber)
		args = append(args, pq.Array(banner.TagIDs))
		argNumber++
	}

	if len(sets) > 0 {
		updateQuery = fmt.Sprintf("update banners set %s where banner_id = $1;", strings.Join(sets, ","))
	}

	res, err := repo.DB.ExecContext(ctx, updateQuery+contentUpdateQuery+tagsUpdateQuery, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (repo *BannersDBRepository) DeleteBanner(ctx context.Context, bannerID int64) (int64, error) {
	res, err := repo.DB.ExecContext(ctx, "delete from banners where banner_id = $1", bannerID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
