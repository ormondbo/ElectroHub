package postgres

import (
	"time"
	"uuid"
)

type dbClient struct {
	id                uuid.UUID `db:"id"`
	name              string    `db:"client_name"`
	surname           string    `db:"client_surname"`
	birthday          time.Time `db:"birthday"`
	gender            string    `db:"gender"`
	registration_date time.Time `db:"registration_date"`
	address_id        uuid.UUID `db:"address_id"`
}

// ======================================================================================================

type dbProduct struct {
	id               uuid.UUID  `db:"id"`
	name             string     `db:"name"`
	category         string     `db:"category"`
	price            float64    `db:"price"`
	available_stock  uint       `db:"available_stock"`
	last_update_date *time.Time `db:"last_update_date"`
	supplier_id      uuid.UUID  `db:"supplier_id"`
	image_id         uuid.UUID  `db:"image_id"`
}

// ======================================================================================================

type dbSupplier struct {
	id           uuid.UUID `db:"id"`
	name         string    `db:"name"`
	address_id   uuid.UUID `db:"address_id"`
	phone_number string    `db:"phone_number"`
}

// ======================================================================================================

type dbImage struct {
	id   uuid.UUID `db:"id"`
	data []byte    `db:"image"`
}

// ======================================================================================================

type dbAddress struct {
	id      uuid.UUID `db:"id"`
	country string    `db:"country"`
	city    string    `db:"city"`
	street  string    `db:"street"`
}
