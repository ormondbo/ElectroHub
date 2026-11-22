package postgres

import (
	"context"
	"electrohub/internal/domain"
	"electrohub/internal/repository"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) repository.ProductRepository {
	return &productRepo{db: db}
}

// ======================================================================================================

func (repo *productRepo) Create(ctx context.Context, dp domain.Product) (uuid.UUID, error) {
	dbp := toDBProduct(dp)
	dbp.id = uuid.New()

	_, err := repo.db.Exec(ctx,
		`INSERT INTO products 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		dbp.id, dbp.name, dbp.category, dbp.price, dbp.available_stock,
		dbp.last_update_date, dbp.supplier_id, dbp.image_id,
	)

	if err != nil {
		return uuid.Nil(), err
	}

	return dbp.id, nil
}

func (repo *productRepo) Update(ctx context.Context, dp domain.Product) error {
	dbp := toDBProduct(dp)

	tag, err := repo.db.Exec(ctx,
		`UPDATE products
			SET name = $1, category = $2, price = $3,
		    	available_stock = $4, last_update_date = $5,
		    	supplier_id = $6, image_id = $7
			WHERE id = $8`,
		dbp.name, dbp.category, dbp.price,
		dbp.available_stock, dbp.last_update_date,
		dbp.supplier_id, dbp.image_id, dbp.id,
	)

	if err == nil && tag.RowsAffected() == 0 {
		err = repository.ErrNotFound
	}

	return err
}

func (repo *productRepo) DecreaseStock(ctx context.Context, id uuid.UUID, amount uint) error {
	tag, err := repo.db.Exec(ctx,
		`UPDATE products
			SET available_stock = available_stock - $2, 
				last_update_date = NOW()
			WHERE id = $1 AND available_stock >= $2`,
		id, amount,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() > 0 {
		return nil
	}

	var exists bool
	if err := repo.db.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 
				FROM products 
					WHERE id = $1
		)`,
		id,
	).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return repository.ErrNotFound
	}

	return repository.ErrInsufficientStock
}

func (repo *productRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	var dbp dbProduct

	err := repo.db.QueryRow(ctx,
		`SELECT *
			FROM products 
			WHERE id = $1`,
		id,
	).Scan(
		&dbp.id, &dbp.name, &dbp.category, &dbp.price, &dbp.available_stock,
		&dbp.last_update_date, &dbp.supplier_id, &dbp.image_id,
	)

	if err != nil {
		return domain.Product{}, err
	}

	return toDomainProduct(dbp), nil
}

func (repo *productRepo) GetAll(ctx context.Context) ([]domain.Product, error) {
	rows, err := repo.db.Query(ctx,
		`SELECT * 
			FROM products
				WHERE available_stock > 0
			ORDER BY id`,
	)

	if err != nil {
		return nil, err
	}

	products := make([]domain.Product, 0, 4)

	defer rows.Close()
	for rows.Next() {
		var dbp dbProduct

		err := rows.Scan(
			&dbp.id, &dbp.name, &dbp.category, &dbp.price, &dbp.available_stock,
			&dbp.last_update_date, &dbp.supplier_id, &dbp.image_id,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, toDomainProduct(dbp))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := repo.db.Exec(ctx,
		`DELETE
		FROM products
			WHERE id = $1`,
		id,
	)

	if err == nil && tag.RowsAffected() == 0 {
		err = repository.ErrNotFound
	}

	return err
}
