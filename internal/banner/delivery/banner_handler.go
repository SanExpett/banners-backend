package delivery

import (
	"context"
	"github.com/SanExpett/banners-backend/internal/banner/usecases"
	"github.com/SanExpett/banners-backend/internal/server/delivery"
	"github.com/SanExpett/banners-backend/pkg/models"
	"github.com/SanExpett/banners-backend/pkg/my_logger"
	"github.com/SanExpett/banners-backend/pkg/utils"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

var _ IBannerService = (*usecases.BannerService)(nil)

type IBannerService interface {
	AddBanner(ctx context.Context, r io.Reader, userID uint64) (uint64, error)
	GetBanner(ctx context.Context, bannerID uint64) (*models.Content, error)
	GetBannersList(ctx context.Context, featureID uint64, tagID uint64, limit uint64,
		offset uint64) ([]*models.Banner, error)
	UpdateBanner(ctx context.Context, r io.Reader, bannerID uint64, userID uint64) error
	DeleteBanner(ctx context.Context, bannerID uint64, userID uint64) error
}

type BannerHandler struct {
	service IBannerService
	logger  *zap.SugaredLogger
}

func NewBannerHandler(bannerService IBannerService) (*BannerHandler, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &BannerHandler{
		service: bannerService,
		logger:  logger,
	}, nil
}

// AddBannerHandler godoc
//
//	@Summary    add banner
//	@Description  add Banner by data
//	@Description Error.status can be:
//	@Description StatusErrBadRequest      = 400
//	@Description  StatusErrInternalServer  = 500
//	@Tags Banner
//
//	@Accept      json
//	@Produce    json
//	@Param      banner  body models.PreBanner true  "Banner data for adding"
//	@Param      token  header string true  "admin token"
//	@Success    200  {object} delivery.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /banner/add [post]
func (b *BannerHandler) AddBannerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	isAdmin, err := delivery.GetIsAdminFromHeader(r)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	if !isAdmin {
		delivery.HandleErr(w, b.logger, delivery.ErrNotAdmin)

		return
	}

	userID, err := delivery.GetUserIDFromHeader(r)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	bannerID, err := b.service.AddBanner(ctx, r.Body, userID)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	delivery.SendOkResponse(w, b.logger, delivery.NewResponseID(bannerID))
	b.logger.Infof("in AddBannerHandler: added banner id= %+v", bannerID)
}

// GetBannerHandler godoc
//
//	@Summary    get banner
//	@Description  get banner by id
//	@Tags Banner
//	@Accept      json
//	@Produce    json
//	@Param      id  query uint64 true  "banner id"
//	@Param      token  header string true  "user token"
//	@Success    200  {object} BannerResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /banner/get [get]
func (b *BannerHandler) GetBannerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	_, err := delivery.GetUserIDFromHeader(r)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	bannerID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	banner, err := b.service.GetBanner(ctx, bannerID)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	delivery.SendOkResponse(w, b.logger, NewBannerResponse(delivery.StatusResponseSuccessful, banner))
	b.logger.Infof("in GetBannerHandler: get Banner: %+v", banner)
}

// DeleteBannerHandler godoc
//
//	@Summary     delete banner
//	@Description  delete banner for author using user id from header\jwt.
//	@Description  This totally removed banner. Recovery will be impossible
//	@Tags Banner
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "banner id"
//	@Param      token  header string true  "admin token"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /banner/delete [delete]
func (b *BannerHandler) DeleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	isAdmin, err := delivery.GetIsAdminFromHeader(r)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	if !isAdmin {
		delivery.HandleErr(w, b.logger, delivery.ErrNotAdmin)

		return
	}

	userID, err := delivery.GetUserIDFromHeader(r)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	bannerIDStr := delivery.GetPathParam(r.URL.Path)
	bannerID, err := strconv.ParseUint(bannerIDStr, 10, 64)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	err = b.service.DeleteBanner(ctx, bannerID, userID)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	delivery.SendOkResponse(w, b.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulDeleteBanner))
	b.logger.Infof("in DeleteBannerHandler: delete Banner id=%d", bannerID)
}

// UpdateBannerHandler godoc
//
//	@Summary    update banner
//	@Description  update banner by data
//	@Description Error.status can be:
//	@Description StatusErrBadRequest      = 400
//	@Description  StatusErrInternalServer  = 500
//	@Tags Banner
//
//	@Accept      json
//	@Produce    json
//	@Param      token  header string true  "admin token"
//	@Param      Banner  body models.PreBanner true  "banner data for updating"
//	@Param      id  path uint64 true  "banner id"
//	@Success    200  {object} delivery.ResponseID
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /banner/update [patch]
func (b *BannerHandler) UpdateBannerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	isAdmin, err := delivery.GetIsAdminFromHeader(r)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	if !isAdmin {
		delivery.HandleErr(w, b.logger, delivery.ErrNotAdmin)

		return
	}

	userID, err := delivery.GetUserIDFromHeader(r)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	bannerIDStr := delivery.GetPathParam(r.URL.Path)
	bannerID, err := strconv.ParseUint(bannerIDStr, 10, 64)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	err = b.service.UpdateBanner(ctx, r.Body, bannerID, userID)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	delivery.SendOkResponse(w, b.logger,
		delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulUpdateBanner))
	b.logger.Infof("in UpdateBannerHandler: added banner id= %+v", bannerID)
}

// GetBannersListHandler godoc
//
//	@Summary    get banners list
//	@Description  get banners list
//	@Tags Banner
//	@Accept      json
//	@Produce    json
//	@Param      feature_id  query uint64 false  "feature_id"
//	@Param      tag_id  query uint64 false  "tag_id"
//	@Param      limit  query uint64 false  "limit Banners"
//	@Param      offset  query uint64 false  "offset of Banners"
//	@Param      token  header string true  "admin token"
//	@Success    200  {object} BannerListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /banner/get_list [get]
func (b *BannerHandler) GetBannersListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	isAdmin, err := delivery.GetIsAdminFromHeader(r)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	if !isAdmin {
		delivery.HandleErr(w, b.logger, delivery.ErrNotAdmin)

		return
	}

	limit, err := utils.ParseUint64FromRequest(r, "limit")
	if err != nil {
		limit = 10
	}

	offset, err := utils.ParseUint64FromRequest(r, "offset")
	if err != nil {
		offset = 0
	}

	featureID, err := utils.ParseUint64FromRequest(r, "feature_id")
	if err != nil {
		featureID = 0
	}

	tagID, err := utils.ParseUint64FromRequest(r, "tag_id")
	if err != nil {
		tagID = 0
	}

	banners, err := b.service.GetBannersList(ctx, featureID, tagID, limit, offset)
	if err != nil {
		delivery.HandleErr(w, b.logger, err)

		return
	}

	delivery.SendOkResponse(w, b.logger, NewBannerListResponse(delivery.StatusResponseSuccessful, banners))
	b.logger.Infof("in GetBannerListHandler: get Banner list: %+v", banners)
}
