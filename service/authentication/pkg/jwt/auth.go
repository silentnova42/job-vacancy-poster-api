package auth

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/silentnova42/job_vacancy_poster/pkg/model"
)

const (
	ExpForRefresh = 7 * 24 * time.Hour
	ExpForAccess  = 15 * time.Minute
)

type AuthService struct {
	refresh string
	access  string
}

func NewAuthService() (*AuthService, error) {
	var a AuthService

	if err := a.FindTokensFromEnv(); err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *AuthService) SetAccessToken(token string) {
	a.access = token
}

func (a *AuthService) SetRefreshToken(token string) {
	a.refresh = token
}

func (a *AuthService) FindTokensFromEnv() error {
	a.refresh = os.Getenv("REFRESH_TOKEN")
	a.access = os.Getenv("ACCESS_TOKEN")

	if strings.TrimSpace(a.refresh) == "" {
		return errors.New("refresh token not found")
	}

	if strings.TrimSpace(a.access) == "" {
		return errors.New("access token not found")
	}

	return nil
}

func (a *AuthService) GenerateAccessToken(customer *model.Customer) (string, error) {
	log.Println(a.access)
	if strings.TrimSpace(a.access) == "" {
		return "", errors.New("we didn't get an access token")
	}
	return generateToken(customer, ExpForAccess, a.access)
}

func (a *AuthService) GenerateRefreshToken(customer *model.Customer) (string, error) {
	log.Println(a.refresh)
	if strings.TrimSpace(a.refresh) == "" {
		return "", errors.New("we didn't get an refresh token")
	}
	return generateToken(customer, ExpForRefresh, a.refresh)
}

func (a *AuthService) Secure(accessToken string) (*model.Customer, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect signature method")
		}
		return []byte(a.access), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email in token not found")
	}

	return &model.Customer{
		Email: email,
	}, nil
}

func (a *AuthService) Refresh(refreshToken string) (*model.Customer, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect signature method")
		}

		return []byte(a.refresh), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	climas, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("impossible to cast token type")
	}

	email, ok := climas["email"].(string)
	if !ok {
		return nil, errors.New("email not found")
	}

	return &model.Customer{
		Email: email,
	}, nil
}

func generateToken(customer *model.Customer, exp time.Duration, key string) (string, error) {
	claims := jwt.MapClaims{
		"email": customer.Email,
		"exp":   exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}
