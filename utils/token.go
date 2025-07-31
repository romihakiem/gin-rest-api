package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"gin-rest-api/config"
	"gin-rest-api/models"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(cfg *config.Config, sub uint) (models.Token, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["sub"] = sub
	claims["exp"] = time.Now().Add(time.Duration(cfg.JWTAccessExpiry) * time.Second).Unix()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token := models.Token{}
	token.AccessToken, err = jwtToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return token, err
	}

	return CreateRefreshToken(cfg, token)
}

func CreateRefreshToken(cfg *config.Config, token models.Token) (models.Token, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, cfg.JWTSecret)

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		fmt.Println(err.Error())
		return token, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return token, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return token, err
	}

	token.RefreshToken = base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(token.AccessToken), nil))
	return token, nil
}

func ValidateToken(cfg *config.Config, accessToken string) (uint, error) {
	jwtToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return 0, err
	}

	result, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		if float64(time.Now().Unix()) > result["exp"].(float64) {
			return 0, errors.New("token expired")
		}

		return uint(result["sub"].(float64)), nil
	}

	return 0, errors.New("token invalid")
}

func ValidateRefreshToken(cfg *config.Config, token models.Token) (uint, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, cfg.JWTSecret)

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return 0, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, err
	}

	data, err := base64.URLEncoding.DecodeString(token.RefreshToken)
	if err != nil {
		return 0, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, err
	}

	if string(plain) != token.AccessToken {
		return 0, errors.New("token invalid")
	}

	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	jwtToken, _, err := parser.ParseUnverified(token.AccessToken, claims)
	if err != nil {
		return 0, err
	}

	result, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		return uint(result["sub"].(float64)), nil
	}

	return 0, errors.New("token invalid")
}
