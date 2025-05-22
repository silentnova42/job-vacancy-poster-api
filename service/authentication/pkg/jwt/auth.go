package auth

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	customer "github.com/silentnova42/job_vacancy_poster/pkg"
)

type Auth struct {
	refresh *string
	access  *string
}

func (a *Auth) FindTokensFromEnv() error {
	*a.refresh = os.Getenv("REFRESH_TOKEN")
	*a.access = os.Getenv("ACCESS_TOKEN")

	if strings.TrimSpace(*a.refresh) == "" {
		return errors.New("refresh token not found")
	}

	if strings.TrimSpace(*a.access) == "" {
		return errors.New("access token not found")
	}

	return nil
}

func (a *Auth) GenerateAccessToken(customer customer.Customer, exp time.Duration) (string, error) {
	return generateToken(customer, exp, *a.access)
}

func (a *Auth) GenerateRefreshToken(customer customer.Customer, exp time.Duration) (string, error) {
	return generateToken(customer, exp, *a.refresh)
}

func (a *Auth) Secure(accessToken string) (*customer.Customer, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect signature method")
		}
		return a.access, nil
	})

	if err != nil || !token.Valid {
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

	password, ok := claims["password"].(string)
	if !ok {
		return nil, errors.New("password in token not found")
	}

	return &customer.Customer{
		Email:    email,
		Password: password,
	}, nil
}

func (a *Auth) Refresh(refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect signature method")
		}

		return a.refresh, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	climas, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token")
	}

	email, ok := climas["email"].(string)
	if !ok {
		return "", "", errors.New("email not found")
	}

	password, ok := climas["password"].(string)
	if !ok {
		return "", "", errors.New("password not found")
	}

	customer := customer.Customer{
		Email:    email,
		Password: password,
	}

	refresh, err := generateToken(customer, 7*24*time.Hour, *a.refresh)
	if err != nil {
		return "", "", err
	}

	access, err := generateToken(customer, 15*time.Minute, *a.access)
	if err != nil {
		return "", "", err
	}

	return refresh, access, nil
}

func generateToken(customer customer.Customer, exp time.Duration, key string) (string, error) {
	claims := jwt.MapClaims{
		"email":    customer.Email,
		"password": customer.Password,
		"exp":      exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}
