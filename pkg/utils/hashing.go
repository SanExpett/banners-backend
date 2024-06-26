package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	myerrors "github.com/SanExpett/banners-backend/pkg/my_errors"
	"github.com/SanExpett/banners-backend/pkg/my_logger"

	"golang.org/x/crypto/argon2"
)

const (
	saltLen = 8

	time    = 1
	memory  = 64 * 1024
	threads = 4
	keyLen  = 32
)

func HashPass(plainPassword string) (string, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}
	salt := make([]byte, saltLen)

	_, err = rand.Read(salt)
	if err != nil {
		logger.Errorln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return hex.EncodeToString(hashPassWithSalt(salt, plainPassword)), nil
}

func hashPassWithSalt(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, time, memory, threads, keyLen)

	return append(salt, hashedPass...)
}

func ComparePassAndHash(passHash []byte, plainPassword string) bool {
	passHashCopy := make([]byte, len(passHash))

	copy(passHashCopy, passHash)

	salt := passHashCopy[0:saltLen]
	userPassHash := hashPassWithSalt(salt[:saltLen], plainPassword)

	return bytes.Equal(userPassHash, passHash)
}
