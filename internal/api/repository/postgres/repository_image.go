package postgres

import (
	"context"
	"electrohub/internal/api/domain"
	"electrohub/internal/api/repository"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type imageRepo struct {
	db *pgxpool.Pool
}

func NewImageRepo(db *pgxpool.Pool) repository.ImageRepository {
	return &imageRepo{db: db}
}

// ======================================================================================================

func (repo *imageRepo) Create(ctx context.Context, di domain.Image) (uuid.UUID, error) {
	dbi := toDBImage(di)
	dbi.id = uuid.New()

	_, err := repo.db.Exec(ctx,
		`INSERT INTO images 
			VALUES ($1, $2)`,
		dbi.id, dbi.data,
	)

	if err != nil {
		return uuid.Nil(), err
	}

	return dbi.id, nil
}

func (repo *imageRepo) Update(ctx context.Context, di domain.Image) error {
	dbi := toDBImage(di)

	tag, err := repo.db.Exec(ctx,
		`UPDATE images
			SET image = $1
				WHERE id = $2`,
		dbi.data, dbi.id,
	)

	if err == nil && tag.RowsAffected() == 0 {
		err = repository.ErrNotFound
	}

	return err
}

func (repo *imageRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Image, error) {
	var dbi dbImage

	err := repo.db.QueryRow(ctx,
		`SELECT * 
		FROM images 
			WHERE id = $1`, id,
	).Scan(
		&dbi.id, &dbi.data,
	)

	if err != nil {
		return domain.Image{}, err
	}

	return toDomainImage(dbi), nil
}

func (repo *imageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := repo.db.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`UPDATE products 
			SET image_id = '00000000-0000-0000-0000-000000000000' 
				WHERE image_id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	tag, err := tx.Exec(ctx,
		`DELETE 
		FROM images
			WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return tx.Commit(ctx)
}
