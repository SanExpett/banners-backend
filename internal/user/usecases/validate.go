package usecases

import (
	"encoding/json"
	"fmt"
	"github.com/SanExpett/banners-backend/pkg/models"
	myerrors "github.com/SanExpett/banners-backend/pkg/my_errors"
	"github.com/SanExpett/banners-backend/pkg/my_logger"
	"github.com/asaskevich/govalidator"
	"io"
)

var (
	ErrWrongCredentials = myerrors.NewError("Некорректный логин (должен быть длиной от 1 до 25 " +
		"символов) или пароль (должен быть не менее 6 символов, содержать цифры, " +
		"строчные и заглавные буквы и специальные символы)")
	ErrDecodeUser = myerrors.NewError("Некорректный json пользователя")
)

func ValidatePreUser(r io.Reader) (*models.PreUser, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	decoder := json.NewDecoder(r)

	preUser := new(models.PreUser)
	if err := decoder.Decode(preUser); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodeUser)
	}

	preUser.Trim()

	_, err = govalidator.ValidateStruct(preUser)
	if err != nil {
		return nil, ErrWrongCredentials
	}

	return preUser, nil
}

func ValidateUserCredentials(login string, password string) (*models.PreUser, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	preUser := new(models.PreUser)

	preUser.Login = login
	preUser.Password = password
	preUser.Trim()
	logger.Infoln(preUser)

	_, err = govalidator.ValidateStruct(preUser)
	if err != nil && (govalidator.ErrorByField(err, "login") != "" ||
		govalidator.ErrorByField(err, "password") != "") {
		logger.Errorln(err)

		return nil, ErrWrongCredentials
	}

	return preUser, nil
}
