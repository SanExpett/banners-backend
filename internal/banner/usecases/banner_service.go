package usecases

import (
	"context"
	"fmt"
	bannerrepo "github.com/SanExpett/banners-backend/internal/banner/repository"
	"github.com/SanExpett/banners-backend/pkg/models"
	myerrors "github.com/SanExpett/banners-backend/pkg/my_errors"
	"github.com/SanExpett/banners-backend/pkg/my_logger"
	"go.uber.org/zap"
	"io"
)

var _ IBannerStorage = (*bannerrepo.BannerStorage)(nil)

type IBannerStorage interface {
	AddBanner(ctx context.Context, preBanner *models.PreBanner, userID uint64) (uint64, error)
	GetBanner(ctx context.Context, bannerID uint64, isAdmin bool) (*models.Content, error)
	GetBannersList(ctx context.Context, featureID uint64, tagID uint64, limit uint64,
		offset uint64) ([]*models.Banner, error)
	UpdateBanner(ctx context.Context, newBanner *models.PreBanner, bannerID uint64, userID uint64) error
	DeleteBanner(ctx context.Context, bannerID uint64, userID uint64) error
}

type BannerService struct {
	storage IBannerStorage
	logger  *zap.SugaredLogger
}

func NewBannerService(bannerStorage IBannerStorage) (*BannerService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &BannerService{storage: bannerStorage, logger: logger}, nil
}

func (b *BannerService) AddBanner(ctx context.Context, r io.Reader, userID uint64) (uint64, error) {
	preBanner, err := ValidatePreBanner(r)
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	bannerID, err := b.storage.AddBanner(ctx, preBanner, userID)
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return bannerID, nil
}

func (b *BannerService) GetBanner(ctx context.Context, bannerID uint64, isAdmin bool) (*models.Content, error) {
	banner, err := b.storage.GetBanner(ctx, bannerID, isAdmin)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	banner.Sanitize()

	return banner, nil
}

func (b *BannerService) DeleteBanner(ctx context.Context, bannerID uint64, userID uint64) error {
	err := b.storage.DeleteBanner(ctx, bannerID, userID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b *BannerService) UpdateBanner(ctx context.Context, r io.Reader, bannerID uint64, userID uint64) error {
	preBanner, err := ValidatePreBanner(r)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	err = b.storage.UpdateBanner(ctx, preBanner, bannerID, userID)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b *BannerService) GetBannersList(ctx context.Context, featureID uint64, tagID uint64, limit uint64,
	offset uint64) ([]*models.Banner, error) {
	banners, err := b.storage.GetBannersList(ctx, featureID, tagID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	for _, banner := range banners {
		banner.Sanitize()
	}

	return banners, nil
}
