package delivery

import (
	"fmt"
	"github.com/SanExpett/banners-backend/pkg/jwt"
	myerrors "github.com/SanExpett/banners-backend/pkg/my_errors"
	"github.com/SanExpett/banners-backend/pkg/my_logger"
	"net/http"
	"strings"
)

func GetIsAdminFromHeader(r *http.Request) (bool, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		logger.Errorln(ErrAuthHeaderNotPresented)

		return false, ErrAuthHeaderNotPresented
	}

	rawJwt := strings.TrimPrefix(authHeader, "Bearer ")

	userPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.Secret)
	if err != nil {
		return false, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return userPayload.IsAdmin, nil
}
