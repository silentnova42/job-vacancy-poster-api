package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/silentnova42/job_vacancy_poster/pkg/structs"
)

func (db *Db) GetAllAvailableVacancy(ctx context.Context) ([]*structs.VacancyGet, error) {
	rows, err := db.client.Query(
		ctx,
		`SELECT id, owner_email, title, description_offer, salary_cents, responses
		FROM public.vacancy;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vacancys := make([]*structs.VacancyGet, 0)
	for rows.Next() {
		var vacancy structs.VacancyGet
		if err = rows.Scan(
			&vacancy.Id,
			&vacancy.OwnerEmail,
			&vacancy.Title,
			&vacancy.DescriptionOffer,
			&vacancy.SalaryCents,
			&vacancy.Responses,
		); err != nil {
			return nil, err
		}
		vacancys = append(vacancys, &vacancy)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vacancys, nil
}

func (db *Db) GetVacancyById(ctx context.Context, id uint) (*structs.VacancyGet, error) {
	var vacancy structs.VacancyGet
	if err := db.client.QueryRow(
		ctx,
		`SELECT id, owner_email, title, description_offer, salary_cents, responses
		FROM public.vacancy
		WHERE id = $1;`,
		id,
	).Scan(
		&vacancy.Id,
		&vacancy.OwnerEmail,
		&vacancy.Title,
		&vacancy.DescriptionOffer,
		&vacancy.SalaryCents,
		&vacancy.Responses,
	); err != nil {
		return nil, err
	}
	return &vacancy, nil
}

func (db *Db) AddVacancy(ctx context.Context, vacancy *structs.VacancyCreate) error {
	_, err := db.client.Exec(
		ctx,
		`INSERT INTO public.vacancy 
		(owner_email, title, description_offer, salary_cents)
		VALUES($1, $2, $3, $4);`,
		vacancy.OwnerEmail,
		vacancy.Title,
		vacancy.DescriptionOffer,
		vacancy.SalaryCents,
	)
	return err
}

func (db *Db) UpdateVacancyById(ctx context.Context, vacancy *structs.VacancyUpdate, id uint) error {
	query, arg, err := buildQuery(vacancy, id)
	if err != nil {
		return err
	}

	_, err = db.client.Exec(ctx, query, arg...)
	return err
}

func buildQuery(vacancy *structs.VacancyUpdate, id uint) (string, []interface{}, error) {
	var (
		query = `UPDATE public.vacancy SET `
		arg   = make([]interface{}, 0)
		parts = make([]string, 0)

		i = 1
	)

	if vacancy.Title != nil {
		parts = append(parts, fmt.Sprintf("title = $%d", i))
		arg = append(arg, *vacancy.Title)
		i++
	}

	if vacancy.DescriptionOffer != nil {
		parts = append(parts, fmt.Sprintf("description_offer = $%d", i))
		arg = append(arg, *vacancy.DescriptionOffer)
		i++
	}

	if vacancy.SalaryCents != nil {
		parts = append(parts, fmt.Sprintf("salary_cents = $%d", i))
		arg = append(arg, *vacancy.SalaryCents)
		i++
	}

	if len(parts) == 0 {
		return "", nil, errors.New("bad request")
	}

	query += strings.Join(parts, ", ")
	query += fmt.Sprintf(" WHERE id = $%d;", i)
	arg = append(arg, id)
	return query, arg, nil
}

func (db *Db) AddResponseById(ctx context.Context, id uint, email string) error {
	tx, err := db.client.Begin(ctx)
	if err != nil {
		return err
	}

	func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if _, err = tx.Exec(
		ctx,
		`UPDATE public.vacancy 
		SET	responses = responses + 1
		WHERE id = $1;`,
		id,
	); err != nil {
		return err
	}

	if _, err = tx.Exec(
		ctx,
		`INSERT INTO public.responses 
		(vacancy_id, email)
		VALUES($1, $2);`,
		id, email,
	); err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func (db *Db) CloseVacancyById(ctx context.Context, id uint) error {
	_, err := db.client.Exec(
		ctx,
		`DELETE FROM public.vacancy
		WHERE id = $1;`,
		id,
	)
	return err
}

func (db *Db) GetResponsesByOwnerId(ctx context.Context, id uint) ([]structs.ResponseGet, error) {
	row, err := db.client.Query(
		ctx,
		`SELECT 
			r.vacancy_id, 
			r.email, 
			v.owner_email
		FROM public.responses AS r
		JOIN public.vacancy AS v ON r.vacancy_id = v.id
		WHERE r.vacancy_id = $1;`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	responses := make([]structs.ResponseGet, 0)
	for row.Next() {
		var response structs.ResponseGet
		if err = row.Scan(&response.VacancyId, &response.Email, &response.OwnerEmail); err != nil {
			return nil, err
		}

		responses = append(responses, response)
	}

	if err = row.Err(); err != nil {
		return nil, err
	}

	return responses, err
}
