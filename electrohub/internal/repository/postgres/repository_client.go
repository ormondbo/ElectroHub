package postgres

import (
	"context"
	"electrohub/internal/domain"
	"electrohub/internal/repository"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type clientRepo struct {
	db *pgxpool.Pool
}

func NewClientRepo(db *pgxpool.Pool) repository.ClientRepository {
	return &clientRepo{db: db}
}

// ======================================================================================================

func (repo *clientRepo) Create(ctx context.Context, dc domain.Client) (uuid.UUID, error) {
	dbc := toDBClient(dc)
	dbc.id = uuid.New()
	_, err := repo.db.Exec(ctx,
		`INSERT INTO clients
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		dbc.id, dbc.name, dbc.surname, dbc.birthday,
		dbc.gender, dbc.registration_date, dbc.address_id,
	)

	if err != nil {
		return uuid.Nil(), err
	}

	return dbc.id, nil
}

func (repo *clientRepo) Update(ctx context.Context, dc domain.Client) error {
	dbc := toDBClient(dc)
	tag, err := repo.db.Exec(ctx,
		`UPDATE clients
			SET client_name = $1, client_surname = $2, birthday = $3,
			gender = $4, registration_date = $5, address_id = $6
				WHERE id = $7`,
		dbc.name, dbc.surname, dbc.birthday, dbc.gender,
		dbc.registration_date, dbc.address_id, dbc.id,
	)

	if err == nil && tag.RowsAffected() == 0 {
		err = repository.ErrNotFound
	}

	return err
}

func (repo *clientRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Client, error) {
	var dbc dbClient
	var dba dbAddress

	err := repo.db.QueryRow(ctx,
		`SELECT c.id, c.client_name, c.client_surname, c.birthday, c.gender, 
			c.registration_date, a.id, a.country, a.city, a.street
				FROM clients c
				JOIN addresses a
					ON c.address_id = a.id 
			WHERE c.id = $1`,
		id,
	).Scan(
		&dbc.id, &dbc.name, &dbc.surname, &dbc.birthday, &dbc.gender,
		&dbc.registration_date, &dba.id, &dba.country, &dba.city, &dba.street,
	)

	if err != nil {
		return domain.Client{}, err
	}

	dc := toDomainClient(dbc)
	dc.Address = toDomainAddress(dba)

	return dc, nil
}

func (repo *clientRepo) GetAll(ctx context.Context, limit uint, offset uint) ([]domain.Client, error) {
	rows, err := repo.db.Query(ctx,
		`SELECT
			c.id, c.client_name, c.client_surname, c.birthday, c.gender,
			c.registration_date, a.id, a.country, a.city, a.street
				FROM clients c 
				JOIN addresses a
					ON c.address_id = a.id
			ORDER BY c.id 
				LIMIT NULLIF($1, 0)
				OFFSET $2`,
		limit, offset,
	)

	if err != nil {
		return nil, err
	}

	clients := make([]domain.Client, 0, 4) //need to test and set capacity by brootforce maybe

	defer rows.Close()
	for rows.Next() {
		var dbc dbClient
		var dba dbAddress

		err := rows.Scan(
			&dbc.id, &dbc.name, &dbc.surname, &dbc.birthday, &dbc.gender,
			&dbc.registration_date, &dba.id, &dba.country, &dba.city, &dba.street,
		)

		if err != nil {
			return nil, err
		}

		dc := toDomainClient(dbc)
		dc.Address = toDomainAddress(dba)
		clients = append(clients, dc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (repo *clientRepo) GetByNameSurname(ctx context.Context, name string, s_name string) ([]domain.Client, error) {
	rows, err := repo.db.Query(ctx,
		`SELECT
			c.id, c.client_name, c.client_surname, c.birthday, c.gender, 
			c.registration_date, a.id, a.country, a.city, a.street
				FROM clients c 
				JOIN addresses a
					ON c.address_id = a.id
			WHERE c.client_name = $1 
				AND c.client_surname = $2
			ORDER BY c.id`,
		name, s_name,
	)

	if err != nil {
		return nil, err
	}

	clients := make([]domain.Client, 0, 1)

	defer rows.Close()
	for rows.Next() {
		var dbc dbClient
		var dba dbAddress

		err := rows.Scan(
			&dbc.id, &dbc.name, &dbc.surname, &dbc.birthday, &dbc.gender,
			&dbc.registration_date, &dba.id, &dba.country, &dba.city, &dba.street,
		)

		if err != nil {
			return nil, err
		}

		dc := toDomainClient(dbc)
		dc.Address = toDomainAddress(dba)
		clients = append(clients, dc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (repo *clientRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := repo.db.Exec(ctx,
		`DELETE 
		FROM clients
			WHERE id = $1`,
		id,
	)

	if err == nil && tag.RowsAffected() == 0 {
		err = repository.ErrNotFound
	}

	return err
}
