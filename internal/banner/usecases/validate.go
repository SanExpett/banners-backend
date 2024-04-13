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
	ErrDecodePreBanner = myerrors.NewError("Некорректный json баннера")
)

func ValidatePreBanner(r io.Reader) (*models.PreBanner, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r)
	preBanner := &models.PreBanner{}
	if err := decoder.Decode(preBanner); err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrDecodePreBanner)
	}

	preBanner.Trim()

	_, err = govalidator.ValidateStruct(preBanner)
	if err != nil {
		logger.Errorln(err)

		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return preBanner, nil
}
