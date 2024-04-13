package repository

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/SanExpett/banners-backend/internal/server/repository"
	"github.com/SanExpett/banners-backend/pkg/models"
	myerrors "github.com/SanExpett/banners-backend/pkg/my_errors"
	"github.com/SanExpett/banners-backend/pkg/my_logger"
	"github.com/SanExpett/banners-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrLoginBusy     = myerrors.NewError("Такой логин уже занят")
	ErrLoginNotExist = myerrors.NewError("Такой логин не существует")
	ErrWrongPassword = myerrors.NewError("Некорректный пароль")

	NameSeqUser = pgx.Identifier{"public", "user_id_seq"} //nolint:gochecknoglobals
)

type UserStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewUserStorage(pool *pgxpool.Pool) (*UserStorage, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	return &UserStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func (u *UserStorage) createUser(ctx context.Context, tx pgx.Tx, preUser *models.PreUser) error {
	var SQLCreateUser string

	var err error

	SQLCreateUser = `INSERT INTO public."user" (login, password) VALUES ($1, $2);`
	_, err = tx.Exec(ctx, SQLCreateUser,
		preUser.Login, preUser.Password)

	if err != nil {
		u.logger.Errorf("in createUser: preUser=%+v err=%+v", preUser, err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (u *UserStorage) AddUser(ctx context.Context, preUser *models.PreUser) (*models.User, error) {
	user := models.User{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		loginBusy, err := u.isLoginBusy(ctx, tx, preUser.Login)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if loginBusy {
			return ErrLoginBusy
		}

		err = u.createUser(ctx, tx, preUser)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		id, err := repository.GetLastValSeq(ctx, tx, NameSeqUser)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		user.ID = id

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	user.Login = preUser.Login
	user.Password = preUser.Password

	return &user, nil
}

func (u *UserStorage) isLoginBusy(ctx context.Context, tx pgx.Tx, login string) (bool, error) {
	SQLIsLoginBusy := `SELECT id FROM public."user" WHERE login=$1;`
	userRow := tx.QueryRow(ctx, SQLIsLoginBusy, login)

	var user string

	if err := userRow.Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		u.logger.Errorln(err)

		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return true, nil
}

func (u *UserStorage) getUserByLogin(ctx context.Context, tx pgx.Tx, login string) (*models.User, error) {
	SQLGetUserByLogin := `SELECT id, login, password FROM public."user" WHERE login=$1;`
	userLine := tx.QueryRow(ctx, SQLGetUserByLogin, login)

	user := models.User{ //nolint:exhaustruct
		Login: login,
	}

	if err := userLine.Scan(&user.ID, &user.Login, &user.Password); err != nil {
		u.logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &user, nil
}

func (u *UserStorage) GetUser(ctx context.Context, login string, password string) (*models.UserWithoutPassword, error) {
	user := &models.User{}                           //nolint:exhaustruct
	userWithoutPass := &models.UserWithoutPassword{} //nolint:exhaustruct

	err := pgx.BeginFunc(ctx, u.pool, func(tx pgx.Tx) error {
		loginBusy, err := u.isLoginBusy(ctx, tx, login)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if !loginBusy {
			return ErrLoginNotExist
		}

		user, err = u.getUserByLogin(ctx, tx, login)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		hashPass, err := hex.DecodeString(user.Password)
		if err != nil {
			return fmt.Errorf(myerrors.ErrTemplate, err)
		}

		if !utils.ComparePassAndHash(hashPass, password) {
			return ErrWrongPassword
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	userWithoutPass.ID = user.ID
	userWithoutPass.Login = user.Login

	return userWithoutPass, nil
}
