package usecases

import (
	"context"
	"fmt"
	userrepo "github.com/SanExpett/banners-backend/internal/user/repository"
	"github.com/SanExpett/banners-backend/pkg/models"
	myerrors "github.com/SanExpett/banners-backend/pkg/my_errors"
	"github.com/SanExpett/banners-backend/pkg/my_logger"
	"github.com/SanExpett/banners-backend/pkg/utils"
	"go.uber.org/zap"
	"io"
)

var _ IUserStorage = (*userrepo.UserStorage)(nil)

type IUserStorage interface {
	AddUser(ctx context.Context, preUser *models.PreUser) (*models.User, error)
	GetUser(ctx context.Context, login string, password string) (*models.UserWithoutPassword, error)
}

type UserService struct {
	storage IUserStorage
	logger  *zap.SugaredLogger
}

func NewUserService(userStorage IUserStorage) (*UserService, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &UserService{storage: userStorage, logger: logger}, nil
}

func (u *UserService) AddUser(ctx context.Context, r io.Reader) (*models.User, error) {
	preUser, err := ValidatePreUser(r)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	preUser.Password, err = utils.HashPass(preUser.Password)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user, err := u.storage.AddUser(ctx, preUser)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return user, nil
}

func (u *UserService) GetUser(ctx context.Context, login string, password string) (*models.UserWithoutPassword, error) {
	preUser, err := ValidateUserCredentials(login, password)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user, err := u.storage.GetUser(ctx, preUser.Login, preUser.Password)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user.Sanitize()

	return user, nil
}
