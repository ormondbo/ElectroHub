package postgres

import (
	"context"
	"electrohub/internal/domain"
	"electrohub/internal/repository"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type addressRepo struct {
	db *pgxpool.Pool
}

func NewAddressRepo(db *pgxpool.Pool) repository.AddressRepository {
	return &addressRepo{db: db}
}

// ======================================================================================================

func (repo *addressRepo) Create(ctx context.Context, da domain.Address) (uuid.UUID, error) {
	dba := toDBAddress(da)
	dba.id = uuid.New()

	_, err := repo.db.Exec(ctx,
		`INSERT INTO addresses 
			VALUES ($1, $2, $3, $4)`,
		dba.id, dba.country, dba.city, dba.street,
	)
	if err != nil {
		return uuid.Nil(), err
	}

	return dba.id, nil
}

func (repo *addressRepo) Update(ctx context.Context, da domain.Address) error {
	dba := toDBAddress(da)

	tag, err := repo.db.Exec(ctx,
		`UPDATE addresses
			SET country = $1, city = $2, street = $3
				WHERE id = $4`,
		dba.country, dba.city, dba.street, dba.id,
	)

	if err == nil && tag.RowsAffected() == 0 {
		err = repository.ErrNotFound
	}

	return err
}

func (repo *addressRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Address, error) {
	var dba dbAddress

	err := repo.db.QueryRow(ctx,
		`SELECT * 
		FROM addresses
			WHERE id = $1`, id,
	).Scan(&dba.id, &dba.country, &dba.city, &dba.street)

	if err != nil {
		return domain.Address{}, err
	}

	return toDomainAddress(dba), nil
}

func (repo *addressRepo) FindByDetails(ctx context.Context, address domain.Address) (uuid.UUID, error) {
	var id uuid.UUID
	err := repo.db.QueryRow(ctx,
		`SELECT id 
		FROM addresses 
			WHERE country = $1 AND city = $2 AND street = $3`,
		address.Country, address.City, address.Street,
	).Scan(&id)

	if err != nil {
		return uuid.Nil(), err
	}

	return id, nil
}

func (repo *addressRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := repo.db.Exec(ctx,
		`DELETE 
		FROM addresses 
			WHERE id = $1`, id,
	)

	if err == nil && tag.RowsAffected() == 0 {
		err = repository.ErrNotFound
	}

	return err
}
