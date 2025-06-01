package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/silentnova42/job_vacancy_poster/service/auth/pkg/model"
)

const (
	ExpForRefresh = 7 * 24 * time.Hour
	ExpForAccess  = 15 * time.Minute
)

type AuthService struct {
	refreshKey string
	accessKey  string
}

func NewAuthService() (*AuthService, error) {
	var a AuthService

	if err := a.FindTokensFromEnv(); err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *AuthService) SetAccessToken(token string) {
	a.accessKey = token
}

func (a *AuthService) SetRefreshToken(token string) {
	a.refreshKey = token
}

func (a *AuthService) FindTokensFromEnv() error {
	a.refreshKey = os.Getenv("REFRESH_TOKEN")
	a.accessKey = os.Getenv("ACCESS_TOKEN")

	if strings.TrimSpace(a.refreshKey) == "" {
		return errors.New("refresh token not found")
	}

	if strings.TrimSpace(a.accessKey) == "" {
		return errors.New("access token not found")
	}

	return nil
}

func (a *AuthService) GenerateAccessToken(customer *model.GetCustomer) (string, error) {
	log.Println(a.accessKey)
	if strings.TrimSpace(a.accessKey) == "" {
		return "", errors.New("we didn't get an access token")
	}
	return generateToken(customer, ExpForAccess, a.accessKey)
}

func (a *AuthService) GenerateRefreshToken(customer *model.GetCustomer) (string, error) {
	log.Println(a.refreshKey)
	if strings.TrimSpace(a.refreshKey) == "" {
		return "", errors.New("we didn't get an refresh token")
	}
	return generateToken(customer, ExpForRefresh, a.refreshKey)
}

func (a *AuthService) Refresh(refreshToken string) (*model.TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect signature method")
		}

		return []byte(a.refreshKey), nil
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

	strId, err := cast(climas, "id")
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		return nil, err
	}

	email, err := cast(climas, "email")
	if err != nil {
		return nil, err
	}

	name, err := cast(climas, "name")
	if err != nil {
		return nil, err
	}

	lastName, err := cast(climas, "last_name")
	if err != nil {
		return nil, err
	}

	customer := model.GetCustomer{
		Id:       uint(id),
		Email:    email,
		Name:     name,
		LastName: lastName,
	}

	access, err := a.GenerateAccessToken(&customer)
	if err != nil {
		return nil, err
	}

	refresh, err := a.GenerateRefreshToken(&customer)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		RefreshToken: refresh,
		AccessToken:  access,
	}, nil
}

func cast(climas jwt.MapClaims, key string) (string, error) {
	value, ok := climas[key]
	if !ok {
		return "", errors.New(key + " not found")
	}

	return fmt.Sprint(value), nil
}

func generateToken(customer *model.GetCustomer, exp time.Duration, key string) (string, error) {
	claims := jwt.MapClaims{
		"id":        customer.Id,
		"email":     customer.Email,
		"name":      customer.Name,
		"last_name": customer.LastName,
		"exp":       exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}
