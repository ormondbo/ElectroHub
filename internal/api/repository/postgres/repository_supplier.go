package postgres

import (
	"context"
	"electrohub/internal/api/domain"
	"electrohub/internal/api/repository"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type supplierRepo struct {
	db *pgxpool.Pool
}

func NewSupplierRepo(db *pgxpool.Pool) repository.SupplierRepository {
	return &supplierRepo{db: db}
}

// ======================================================================================================

func (repo *supplierRepo) Create(ctx context.Context, ds domain.Supplier) (uuid.UUID, error) {
	dbs := toDBSupplier(ds)
	dbs.id = uuid.New()

	_, err := repo.db.Exec(ctx,
		`INSERT INTO suppliedbs
			VALUES ($1, $2, $3, $4)`,
		dbs.id, dbs.name, dbs.address_id, dbs.phone_number,
	)

	if err != nil {
		return uuid.Nil(), err
	}

	return dbs.id, nil
}

func (repo *supplierRepo) Update(ctx context.Context, ds domain.Supplier) error {
	dbs := toDBSupplier(ds)

	tag, err := repo.db.Exec(ctx,
		`UPDATE suppliedbs
			SET name = $1, address_id = $2, phone_number = $3
				WHERE id = $4`,
		dbs.name, dbs.address_id, dbs.phone_number, dbs.id,
	)

	if err == nil && tag.RowsAffected() == 0 {
		err = repository.ErrNotFound
	}

	return err
}

func (repo *supplierRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Supplier, error) {
	var dbs dbSupplier
	var dba dbAddress

	err := repo.db.QueryRow(ctx,
		`SELECT s.id, s.name, s.address_id, s.phone_number, a.id, a.country, a.city, a.street
			FROM suppliedbs s
			JOIN addresses a
				ON s.address_id = a.id
			WHERE s.id = $1`, id,
	).Scan(
		&dbs.id, &dbs.name, &dbs.address_id, &dbs.phone_number,
		&dba.id, &dba.country, &dba.city, &dba.street,
	)

	if err != nil {
		return domain.Supplier{}, err
	}

	supplier := toDomainSupplier(dbs)
	supplier.Address = toDomainAddress(dba)

	return supplier, nil
}

func (repo *supplierRepo) GetAll(ctx context.Context) ([]domain.Supplier, error) {
	rows, err := repo.db.Query(ctx,
		`SELECT
			s.id, s.name, s.address_id, s.phone_number, a.id, a.country, a.city, a.street
				FROM suppliedbs s
				JOIN addresses a
					ON s.address_id = a.id
			ORDER BY s.id`,
	)

	if err != nil {
		return nil, err
	}

	suppliers := make([]domain.Supplier, 0, 4) //need to test and set capacity by brootforce maybe

	defer rows.Close()
	for rows.Next() {
		var dbs dbSupplier
		var dba dbAddress

		err := rows.Scan(
			&dbs.id, &dbs.name, &dbs.address_id, &dbs.phone_number,
			&dba.id, &dba.country, &dba.city, &dba.street,
		)

		if err != nil {
			return nil, err
		}

		supplier := toDomainSupplier(dbs)
		supplier.Address = toDomainAddress(dba)
		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suppliers, nil
}

func (repo *supplierRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := repo.db.Exec(ctx,
		`DELETE
		FROM suppliedbs
			WHERE id = $1`,
		id,
	)

	if err == nil && tag.RowsAffected() == 0 {
		err = repository.ErrNotFound
	}

	return err
}
