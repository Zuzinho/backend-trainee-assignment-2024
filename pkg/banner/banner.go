package banner

import (
	"avito_hr/pkg/user"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Content map[string]any

func (content *Content) Value() (driver.Value, error) {
	return json.Marshal(content)
}

func (content *Content) Scan(value any) error {
	buf, ok := value.([]byte)
	if !ok {
		return IncorrectTypeFromDBErr
	}

	return json.Unmarshal(buf, content)
}

type Banner struct {
	BannerID  int64     `json:"banner_id"`
	FeatureID int64     `json:"feature_id"`
	TagIDs    []int64   `json:"tag_ids"`
	Content   *Content  `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBanner struct {
	FeatureID int64    `json:"feature_id"`
	TagIDs    []int64  `json:"tag_ids"`
	Content   *Content `json:"content"`
	IsActive  bool     `json:"is_active"`
}

type UpdateBanner struct {
	FeatureID *int64   `json:"feature_id"`
	TagIDs    *[]int64 `json:"tag_ids"`
	Content   *Content `json:"content"`
	IsActive  *bool    `json:"is_active"`
}

type Banners []*Banner

func (banners *Banners) Append(banner *Banner) {
	*banners = append(*banners, banner)
}

type BannersRepository interface {
	GetContent(ctx context.Context, tagID, featureID sql.NullInt64, role user.Role) (*Content, error)
	GetBanners(ctx context.Context, tagID, featureID, limit, offset sql.NullInt64) (*Banners, error)
	CreateBanner(ctx context.Context, banner *CreateBanner) (int64, error)
	UpdateBanner(ctx context.Context, bannerID int64, banner *UpdateBanner) (int64, error)
	DeleteBanner(ctx context.Context, bannerID int64) (int64, error)
}

type BannersTempRepository interface {
	UpdateBanners(ticker *time.Ticker)
	GetContent(tagID, featureID sql.NullInt64, role user.Role) *Content
}
