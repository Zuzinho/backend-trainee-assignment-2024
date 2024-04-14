package banner

import (
	"avito_hr/pkg/user"
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"sort"
	"sync"
	"time"
)

type BannersTempMemoryRepository struct {
	mu          *sync.RWMutex
	data        map[int64]*Banner
	BannersRepo BannersRepository
}

func NewBannersTempMemoryRepository(repository BannersRepository) *BannersTempMemoryRepository {
	return &BannersTempMemoryRepository{
		mu:          &sync.RWMutex{},
		data:        make(map[int64]*Banner),
		BannersRepo: repository,
	}
}

func (repo *BannersTempMemoryRepository) UpdateBanners(ticker *time.Ticker) {
	null := sql.NullInt64{}

	for {
		select {
		case <-ticker.C:
			log.Info("Trying update temp memory repository...")

			banners, err := repo.BannersRepo.GetBanners(context.Background(),
				null, null, null, null)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err.Error(),
				}).Error("Caught error while updating temp repository, try again in 5 minutes")
				return
			}

			repo.mu.Lock()
			repo.data = make(map[int64]*Banner)
			for _, banner := range *banners {
				repo.data[banner.BannerID] = banner
			}
			repo.mu.Unlock()
		}
	}
}

func (repo *BannersTempMemoryRepository) GetContent(tagID, featureID sql.NullInt64, role user.Role) *Content {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, v := range repo.data {
		_, found := sort.Find(len(v.TagIDs), func(i int) int {
			return int(tagID.Int64 - v.TagIDs[i])
		})

		if v.FeatureID == featureID.Int64 && found {
			if role == user.RoleUser && !v.IsActive {
				return nil
			}

			return v.Content
		}
	}

	return nil
}
