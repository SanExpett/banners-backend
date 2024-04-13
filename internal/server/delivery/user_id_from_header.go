package delivery

import (
	"fmt"
	"github.com/SanExpett/banners-backend/pkg/jwt"
	myerrors "github.com/SanExpett/banners-backend/pkg/my_errors"
	"github.com/SanExpett/banners-backend/pkg/my_logger"
	"net/http"
	"strings"
)

func GetUserIDFromHeader(r *http.Request) (uint64, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		logger.Errorln(ErrAuthHeaderNotPresented)

		return 0, ErrAuthHeaderNotPresented
	}

	rawJwt := strings.TrimPrefix(authHeader, "Bearer ")

	userPayload, err := jwt.NewUserJwtPayload(rawJwt, jwt.Secret)
	if err != nil {
		return 0, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return userPayload.UserID, nil
}
