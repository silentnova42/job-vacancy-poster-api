package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/silentnova42/job_vacancy_poster/service/vacancy/pkg/model"
)

func (db *Db) GetAllAvailableVacancy(ctx context.Context) ([]*model.VacancyGet, error) {
	rows, err := db.client.Query(
		ctx,
		`SELECT 
			id
			, owner_email
			, title
			, description_offer
			, salary_cents
		FROM public.vacancies;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vacancys := make([]*model.VacancyGet, 0)
	for rows.Next() {
		var vacancy model.VacancyGet
		if err = rows.Scan(
			&vacancy.Id,
			&vacancy.OwnerEmail,
			&vacancy.Title,
			&vacancy.DescriptionOffer,
			&vacancy.SalaryCents,
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

func (db *Db) GetVacancyById(ctx context.Context, vacancyId uint) (*model.VacancyGetWithResponses, error) {
	var vacancy model.VacancyGetWithResponses

	count, err := db.GetCountResponsesById(ctx, vacancyId)
	if err != nil {
		return nil, err
	}
	vacancy.Responses = count

	if err := db.client.QueryRow(
		ctx,
		`SELECT 
			id
			, owner_email
			, title
			, description_offer
			, salary_cents
		FROM public.vacancies
		WHERE id = $1;`,
		vacancyId,
	).Scan(
		&vacancy.Id,
		&vacancy.OwnerEmail,
		&vacancy.Title,
		&vacancy.DescriptionOffer,
		&vacancy.SalaryCents,
	); err != nil {
		return nil, err
	}
	return &vacancy, nil
}

func (db *Db) AddVacancy(ctx context.Context, vacancy *model.VacancyCreate, email string) error {
	_, err := db.client.Exec(
		ctx,
		`INSERT INTO public.vacancies 
			( owner_email
			, title
			, description_offer
			, salary_cents )
		VALUES($1, $2, $3, $4);`,
		email,
		vacancy.Title,
		vacancy.DescriptionOffer,
		vacancy.SalaryCents,
	)
	return err
}

func (db *Db) UpdateVacancyByIdAndEmail(ctx context.Context, vacancy *model.VacancyUpdate, vacancyId uint, email string) error {
	query, arg, err := buildQuery(vacancy, vacancyId, email)
	if err != nil {
		return err
	}

	_, err = db.client.Exec(ctx, query, arg...)
	return err
}

func buildQuery(vacancy *model.VacancyUpdate, vacancyId uint, email string) (string, []interface{}, error) {
	var (
		query = `UPDATE public.vacancies SET `
		arg   = make([]interface{}, 0)
		parts = make([]string, 0)
		index = 1
	)

	if vacancy.Title != nil {
		parts = append(parts, fmt.Sprintf("title = $%d", index))
		arg = append(arg, *vacancy.Title)
		index++
	}

	if vacancy.DescriptionOffer != nil {
		parts = append(parts, fmt.Sprintf("description_offer = $%d", index))
		arg = append(arg, *vacancy.DescriptionOffer)
		index++
	}

	if vacancy.SalaryCents != nil {
		parts = append(parts, fmt.Sprintf("salary_cents = $%d", index))
		arg = append(arg, *vacancy.SalaryCents)
		index++
	}

	if len(parts) == 0 {
		return "", nil, errors.New("bad request")
	}

	query += strings.Join(parts, ", ")
	query += fmt.Sprintf(" WHERE id = $%d ", index)
	index++
	query += fmt.Sprintf("AND email = $%d;", index)
	arg = append(arg, vacancyId, email)
	return query, arg, nil
}

func (db *Db) CloseVacancyByIdAndEmail(ctx context.Context, id uint, email string) error {
	_, err := db.client.Exec(
		ctx,
		`DELETE FROM public.vacancies
		WHERE id = $1 AND email = $2;`,
		id, email,
	)
	return err
}

func (db *Db) GetResponsesByVacancyId(ctx context.Context, vacancyId uint) ([]model.ResponseGet, error) {
	row, err := db.client.Query(
		ctx,
		`SELECT 
			r.vacancy_id, 
			r.email, 
			v.owner_email
		FROM public.responses AS r
		JOIN public.vacancies AS v 
		ON r.vacancy_id = v.id
		WHERE r.vacancy_id = $1;`,
		vacancyId,
	)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	responses := make([]model.ResponseGet, 0)
	for row.Next() {
		var response model.ResponseGet
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

func (db *Db) GetCountResponsesById(ctx context.Context, vacancyId uint) (int, error) {
	var countResponses int

	err := db.client.QueryRow(
		ctx,
		`SELECT COUNT(*) 
		FROM public.responses
		WHERE vacancy_id = $1;`,
		vacancyId,
	).Scan(&countResponses)

	return countResponses, err
}

func (db *Db) AddResponseByIdAndEmail(ctx context.Context, vacancyId uint, email string) error {
	var id uint

	err := db.client.QueryRow(ctx,
		`SELECT id 
		FROM public.responses 
		WHERE vacancy_id = $1 AND email = $2`,
		vacancyId, email,
	).Scan(&id)
	if err == nil {
		return errors.New("you have already applied")
	}

	if errors.Is(err, pgx.ErrNoRows) {
		_, err = db.client.Exec(
			ctx,
			`INSERT INTO public.responses 
			( vacancy_id
			, email )
		VALUES($1, $2);`,
			vacancyId, email,
		)
		return err
	}

	return err
}

func (db *Db) DeleteResponseByIdAndEmail(ctx context.Context, vacancyId uint, email string) error {
	_, err := db.client.Exec(
		ctx,
		`DELETE 
		FROM public.responses 
		WHERE vacancy_id = $1 AND email = $2`,
		vacancyId, email,
	)
	return err
}
