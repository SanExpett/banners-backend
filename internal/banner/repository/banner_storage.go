package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/SanExpett/banners-backend/internal/server/repository"
	"github.com/SanExpett/banners-backend/pkg/models"
	myerrors "github.com/SanExpett/banners-backend/pkg/my_errors"
	"github.com/SanExpett/banners-backend/pkg/my_logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrBannerNotFound       = myerrors.NewError("Этот баннер не найден")
	ErrNoAffectedBannerRows = myerrors.NewError("Не получилось обновить данные баннера")

	NameSeqBanner = pgx.Identifier{"public", "banner_id_seq"} //nolint:gochecknoglobals
)

type BannerStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewBannerStorage(pool *pgxpool.Pool) (*BannerStorage, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &BannerStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func (b *BannerStorage) createBanner(ctx context.Context, tx pgx.Tx, preBanner *models.PreBanner,
	userID uint64) error {
	var SQLCreateBanner string

	var err error

	SQLCreateBanner = `INSERT INTO public."banner" (author_id, feature_id, 
                             title, text, url, is_active) VALUES ($1, $2, $3, $4, $5, $6);`
	_, err = tx.Exec(ctx, SQLCreateBanner, userID, preBanner.FeatureID,
		preBanner.Content.Title, preBanner.Content.Text, preBanner.Content.URL, preBanner.IsActive)

	if err != nil {
		b.logger.Errorf("in createBanner: preBanner%+v err=%+v", preBanner, err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b *BannerStorage) addTag(ctx context.Context, tx pgx.Tx, tagID, bannerID uint64) error {
	var SQLAddTag string

	var err error

	SQLAddTag = `INSERT INTO public."banner_tag" (banner_id, tag_id) VALUES ($1, $2);`
	_, err = tx.Exec(ctx, SQLAddTag, bannerID, tagID)

	if err != nil {
		b.logger.Errorf("in addTag: tagID=%d bannerID%d", tagID, bannerID)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b *BannerStorage) AddBanner(ctx context.Context, preBanner *models.PreBanner, userID uint64) (uint64, error) {
	var bannerID uint64

	err := pgx.BeginFunc(ctx, b.pool, func(tx pgx.Tx) error {
		err := b.createBanner(ctx, tx, preBanner, userID)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		id, err := repository.GetLastValSeq(ctx, tx, NameSeqBanner)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		for _, tagID := range preBanner.TagIDs {
			err = b.addTag(ctx, tx, tagID, bannerID)
			if err != nil {
				return fmt.Errorf(myerrors.ErrTemplate, err)
			}
		}

		bannerID = id

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return bannerID, nil
}

func (b *BannerStorage) selectBannerContentByID(ctx context.Context,
	tx pgx.Tx, bannerID uint64,
) (*models.Content, error) {
	SQLSelectBanner := `SELECT title, text, url FROM public."banner" WHERE id=$1`
	bannerContent := &models.Content{} //nolint:exhaustruct

	BannerRow := tx.QueryRow(ctx, SQLSelectBanner, bannerID)
	if err := BannerRow.Scan(&bannerContent.Title, &bannerContent.Text, &bannerContent.URL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(myerrors.ErrTemplate, ErrBannerNotFound)
		}

		b.logger.Errorf("error with bannerId=%d: %+v", bannerID, err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return bannerContent, nil
}

func (b *BannerStorage) GetBanner(ctx context.Context, bannerID uint64) (*models.Content, error) {
	var bannerContent *models.Content

	err := pgx.BeginFunc(ctx, b.pool, func(tx pgx.Tx) error {
		bannerContentInner, err := b.selectBannerContentByID(ctx, tx, bannerID)
		if err != nil {
			return err
		}

		bannerContent = bannerContentInner

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return bannerContent, nil
}

func (b *BannerStorage) deleteBanner(ctx context.Context, tx pgx.Tx, bannerID uint64, userID uint64) error {
	SQLDeleteBanner := `DELETE FROM public."banner" WHERE id=$1 AND author_id=$2`

	result, err := tx.Exec(ctx, SQLDeleteBanner, bannerID, userID)
	if err != nil {
		b.logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(myerrors.ErrTemplate, ErrNoAffectedBannerRows)
	}

	return nil
}

func (b *BannerStorage) DeleteBanner(ctx context.Context, bannerID uint64, userID uint64) error {
	err := pgx.BeginFunc(ctx, b.pool, func(tx pgx.Tx) error {
		err := b.deleteBanner(ctx, tx, bannerID, userID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		b.logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b *BannerStorage) updateBanner(ctx context.Context, tx pgx.Tx, preBanner *models.PreBanner,
	bannerID uint64, userID uint64) error {
	var SQLCreateBanner string

	var err error

	SQLCreateBanner = `UPDATE public."banner" SET feature_id = $1, title = $2, text = $3, url = $4, is_active = $5 
                             WHERE author_id=$7 AND id=$8;`
	_, err = tx.Exec(ctx, SQLCreateBanner, preBanner.FeatureID,
		preBanner.Content.Title, preBanner.Content.Text, preBanner.Content.URL, preBanner.IsActive, userID, bannerID)

	if err != nil {
		b.logger.Errorf("in updateBanner: preBanner%+v err=%+v", preBanner, err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b *BannerStorage) UpdateBanner(ctx context.Context, newBanner *models.PreBanner, bannerID uint64,
	userID uint64) error {
	err := pgx.BeginFunc(ctx, b.pool, func(tx pgx.Tx) error {
		err := b.updateBanner(ctx, tx, newBanner, userID, bannerID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		b.logger.Errorln(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (b *BannerStorage) selectTagsIDsByBannerID(ctx context.Context, tx pgx.Tx,
	bannerID uint64) ([]uint64, error) {
	SQLSelectTagssIDsByBannerID :=
		`SELECT tag_id
		FROM public."banner_tag" 
		WHERE banner_id = $1`

	tagsIDsByBannerIDRows, err := tx.Query(ctx, SQLSelectTagssIDsByBannerID, bannerID)
	if err != nil {
		b.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	var curTagID uint64
	var slTagIDs []uint64

	_, err = pgx.ForEachRow(tagsIDsByBannerIDRows, []any{
		&curTagID,
	}, func() error {
		slTagIDs = append(slTagIDs, curTagID)

		return nil
	})
	if err != nil {
		b.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slTagIDs, nil
}

func (b *BannerStorage) selectBannersInFeedWithWhereLimitOffset(ctx context.Context, tx pgx.Tx,
	featureID uint64, tagID uint64, limit uint64, offset uint64) ([]*models.Banner, error) {
	query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("id, feature_id, " +
		"title, text, url, is_active, created_at, updated_at").From(`public."banner"`)

	if featureID != 0 || tagID != 0 {
		if featureID != 0 {
			query = query.Where(squirrel.Eq{"feature_id": featureID})
		}
		if tagID != 0 {
			query = query.Join(`public."banner_tag" bt ON public."banner".id = bt.banner_id`).
				Join(`public."tag" t ON bt.tag_id = t.id`).
				Where(squirrel.Eq{"t.id": tagID})
		}
	}

	query = query.Limit(limit).Offset(offset)

	SQLQuery, args, err := query.ToSql()
	if err != nil {
		b.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	rowsBanners, err := tx.Query(ctx, SQLQuery, args...)
	if err != nil {
		b.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	curBanner := new(models.Banner)

	var slBanner []*models.Banner

	_, err = pgx.ForEachRow(rowsBanners, []any{
		&curBanner.BannerID, &curBanner.FeatureID,
		&curBanner.Content.Title, &curBanner.Content.Text, &curBanner.Content.URL,
		&curBanner.IsActive, &curBanner.CreatedAt, &curBanner.UpdatedAt,
	}, func() error {
		slBanner = append(slBanner, &models.Banner{
			BannerID:  curBanner.BannerID,
			FeatureID: curBanner.FeatureID,
			Content:   curBanner.Content,
			IsActive:  curBanner.IsActive,
			CreatedAt: curBanner.CreatedAt,
			UpdatedAt: curBanner.UpdatedAt,
		})

		return nil
	})
	if err != nil {
		b.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slBanner, nil
}

func (b *BannerStorage) GetBannersList(ctx context.Context, featureID uint64, tagID uint64, limit uint64,
	offset uint64) ([]*models.Banner, error) {
	var slBanners []*models.Banner

	err := pgx.BeginFunc(ctx, b.pool, func(tx pgx.Tx) error {
		slBannersInner, err := b.selectBannersInFeedWithWhereLimitOffset(ctx, tx, featureID, tagID, limit, offset)
		if err != nil {
			return err
		}

		for _, banner := range slBannersInner {
			banner.TagIDs, err = b.selectTagsIDsByBannerID(ctx, tx, banner.BannerID)
			if err != nil {
				return err
			}
		}

		slBanners = slBannersInner

		return nil
	})
	if err != nil {
		b.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return slBanners, nil
}
