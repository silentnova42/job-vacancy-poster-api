package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/silentnova42/job_vacancy_poster/service/profile/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

func (db *Db) GetCustomerByEmailAndPassword(ctx context.Context, loginRequest *model.LoginRequest) (*model.GetPrivateCustomer, error) {
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
		loginRequest.Email,
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

	if err = comperePasswordHash([]byte(customer.Password), []byte(loginRequest.Password)); err != nil {
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

func (db *Db) AddCustomer(ctx context.Context, newCustomer *model.CreateCustomer) error {
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

func (db *Db) UpdateCustomer(ctx context.Context, updateCustomer *model.UpdateCustomer, email string) error {
	customer, err := db.GetCustomerByEmail(ctx, email)
	if err != nil {
		return err
	}

	query, args, err := buildQuery(updateCustomer, customer)
	if err != nil {
		return err
	}

	_, err = db.client.Exec(ctx, query, args...)
	return err
}

func buildQuery(updateCustomer *model.UpdateCustomer, customer *model.GetPublicCustomer) (string, []interface{}, error) {
	var (
		query = "UPDATE public.profiles SET "
		parts = make([]string, 0)
		args  = make([]interface{}, 0)
		index = 1
	)

	if updateCustomer.NewName != nil && *updateCustomer.NewName != customer.Name {
		parts = append(parts, fmt.Sprintf("name = $%v", index))
		args = append(args, updateCustomer.NewName)
		index++
	}

	if updateCustomer.NewLastName != nil && *updateCustomer.NewLastName != customer.LastName {
		parts = append(parts, fmt.Sprintf("last_name = $%v", index))
		args = append(args, updateCustomer.NewLastName)
		index++
	}

	if updateCustomer.NewResume != nil && *updateCustomer.NewResume != customer.Resume {
		parts = append(parts, fmt.Sprintf("resume = $%v", index))
		args = append(args, updateCustomer.NewResume)
		index++
	}

	if len(parts) == 0 {
		return "", nil, errors.New("there is nothing to change")
	}

	query += strings.Join(parts, ", ")
	query += fmt.Sprintf(" WHERE email = $%d;", index)
	args = append(args, customer.Email)
	return query, args, nil
}

func (db *Db) UpdatePassword(ctx context.Context, passwordUpdate *model.PasswordUpdateRequest, email string) error {
	var (
		passwordHash string
		err          error
	)

	if err = db.client.QueryRow(
		ctx,
		`
		SELECT password 
		FROM public.profiles
		WHERE email = $1;
		`,
		email,
	).Scan(&passwordHash); err != nil {
		return err
	}

	if err = comperePasswordHash([]byte(passwordHash), []byte(passwordUpdate.OldPassword)); err != nil {
		return err
	}

	newPasswordHash, err := getPasswordHash(passwordUpdate.NewPassword)
	if err != nil {
		return err
	}

	if _, err = db.client.Exec(
		ctx,
		`
		UPDATE public.profiles 
		SET password = $1 
		WHERE email = $2;
		`,
		newPasswordHash,
		email,
	); err != nil {
		return err
	}

	return nil
}

func (db *Db) DeleteCustomerByEmailAndPassword(ctx context.Context, credentials *model.PasswordPayload, email string) error {
	var (
		passwordHash string
		err          error
	)

	if err = db.client.QueryRow(
		ctx,
		`
		SELECT password 
		FROM public.profiles
		WHERE email = $1;
		`,
		email,
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
		email, passwordHash,
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
