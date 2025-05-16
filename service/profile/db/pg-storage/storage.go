package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/silentnova42/job_vacancy_poster/pkg/structs"
)

func (db *Db) GetProfileByEmailAndPassword(ctx context.Context, checkCustomer structs.CheckCustomer) (*structs.GetCustomer, error) {
	var customer structs.GetCustomer

	if err := db.client.QueryRow(
		ctx,
		`SELECT 
			email
			, name
			, last_name
			, password
		FROM public.profiles
		WHERE email=$1`,
		checkCustomer.Email,
	).Scan(
		&customer.Email,
		&customer.Name,
		&customer.LastName,
		&customer.Password,
	); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (db *Db) AddProfile(ctx context.Context, customer structs.CreateCustomer) error {
	_, err := db.client.Exec(
		ctx,
		`INSERT INTO public.profiles 
			( email
			, name
			, last_name 
			, password )
		VALUES($1, $2, $3, $4)`,
		customer.Email,
		customer.Name,
		customer.LastName,
		customer.Password,
	)
	return err
}

func (db *Db) UpdateProfile(ctx context.Context, updateCustomer structs.UpdateCustomer, getCustomer structs.GetCustomer) error {
	query, args, err := buildQuery(updateCustomer, getCustomer)
	if err != nil {
		return err
	}

	_, err = db.client.Exec(ctx, query, args...)
	return err
}

func buildQuery(updateCustomer structs.UpdateCustomer, getCustomer structs.GetCustomer) (string, []interface{}, error) {
	var (
		query = "UPDATE public.profiles SET "
		parts = make([]string, 0)
		args  = make([]interface{}, 0)
		index = 1
	)

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

	if len(parts) == 0 {
		return "", nil, errors.New("there is nothing to change")
	}

	query += strings.Join(parts, ", ")
	return query, args, nil
}
