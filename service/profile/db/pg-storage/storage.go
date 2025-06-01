package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/silentnova42/job_vacancy_poster/service/profile/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

func (db *Db) GetProfileByEmailAndPassword(ctx context.Context, credentials *model.Credentials) (*model.GetPrivateCustomer, error) {
	var (
		customer model.GetPrivateCustomer
		err      error
	)

	if err = db.client.QueryRow(
		ctx,
		`SELECT
			id
			, email
			, name
			, last_name
			, resume
			, password
		FROM public.profiles
		WHERE email=$1`,
		credentials.Email,
	).Scan(
		&customer.Id,
		&customer.Email,
		&customer.Name,
		&customer.LastName,
		&customer.Resume,
		&customer.Password,
	); err != nil {
		return nil, err
	}

	if err = comperePasswordHash([]byte(customer.Password), []byte(credentials.Password)); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (db *Db) GetCustomerByEmail(ctx context.Context, email string) (*model.GetPublicCustomer, error) {
	var customer model.GetPublicCustomer

	if err := db.client.QueryRow(
		ctx,
		`
		SELECT
			email
			, name
			, last_name
			, resume
		FROM public.profiles
		WHERE email = $1 
		`,
		email,
	).Scan(
		&customer.Email,
		&customer.Name,
		&customer.LastName,
		&customer.Resume,
	); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (db *Db) AddProfile(ctx context.Context, newCustomer *model.CreateCustomer) error {
	hash, err := getPasswordHash(newCustomer.Password)
	if err != nil {
		return err
	}

	_, err = db.client.Exec(
		ctx,
		`INSERT INTO public.profiles 
			( email
			, name
			, last_name
			, resume
			, password )
		VALUES($1, $2, $3, $4, $5)`,
		newCustomer.Email,
		newCustomer.Name,
		newCustomer.LastName,
		newCustomer.Resume,
		hash,
	)
	return err
}

func (db *Db) UpdateProfile(ctx context.Context, updateCustomer *model.UpdateCustomer) error {
	getCustomer, err := db.GetProfileByEmailAndPassword(ctx, &updateCustomer.Credentials)
	if err != nil {
		return err
	}

	query, args, err := buildQuery(updateCustomer, getCustomer)
	if err != nil {
		return err
	}

	_, err = db.client.Exec(ctx, query, args...)
	return err
}

func buildQuery(updateCustomer *model.UpdateCustomer, getCustomer *model.GetPrivateCustomer) (string, []interface{}, error) {
	var (
		query = "UPDATE public.profiles SET "
		parts = make([]string, 0)
		args  = make([]interface{}, 0)
		index = 1
	)

	if err := comperePasswordHash([]byte(getCustomer.Password), []byte(updateCustomer.Credentials.Password)); err != nil {
		return "", nil, err
	}

	if updateCustomer.Email != nil && *updateCustomer.Email != getCustomer.Email {
		parts = append(parts, fmt.Sprintf("email = $%v", index))
		args = append(args, updateCustomer.Email)
		index++
	}

	if updateCustomer.Name != nil && *updateCustomer.Name != getCustomer.Name {
		parts = append(parts, fmt.Sprintf("name = $%v", index))
		args = append(args, updateCustomer.Name)
		index++
	}

	if updateCustomer.LastName != nil && *updateCustomer.LastName != getCustomer.LastName {
		parts = append(parts, fmt.Sprintf("last_name = $%v", index))
		args = append(args, updateCustomer.LastName)
		index++
	}

	if updateCustomer.Resume != nil && *updateCustomer.Resume != getCustomer.Resume {
		parts = append(parts, fmt.Sprintf("resume = $%v", index))
		args = append(args, updateCustomer.Resume)
		index++
	}

	log.Println(*updateCustomer.Password)
	log.Println(getCustomer.Password)

	if updateCustomer.Password != nil && *updateCustomer.Password != updateCustomer.Credentials.Password {
		parts = append(parts, fmt.Sprintf("password = $%v", index))

		hash, err := getPasswordHash(*updateCustomer.Password)
		if err != nil {
			return "", nil, err
		}

		args = append(args, hash)
		index++
	}

	if len(parts) == 0 {
		return "", nil, errors.New("there is nothing to change")
	}

	query += strings.Join(parts, ", ")
	return query, args, nil
}

func (db *Db) DeleteProfileByEmailAndPassword(ctx context.Context, credentials *model.Credentials) error {
	var (
		passwordHash string
		err          error
	)

	if err = db.client.QueryRow(
		ctx,
		`
		SELECT password FROM public.profiles
		WHERE email = $1
		`,
		credentials.Email,
	).Scan(&passwordHash); err != nil {
		return err
	}

	if err = comperePasswordHash([]byte(passwordHash), []byte(credentials.Password)); err != nil {
		return err
	}

	_, err = db.client.Exec(
		ctx,
		`DELETE FROM public.profiles
		WHERE email = $1 AND password = $2`,
		credentials.Email, passwordHash,
	)

	return err
}

func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func comperePasswordHash(hash []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
